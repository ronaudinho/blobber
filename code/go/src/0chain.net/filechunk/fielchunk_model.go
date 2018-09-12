package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"runtime"
	"sort"
	"sync"

	"os"

	// . "0chain.net/logging"
	"github.com/klauspost/reedsolomon"
)

//FileInfo struct
type FileInfo struct {
	DataShards int
	ParShards  int
	File       string
	OutDir     string
}

//Maininfo is for main response
type Maininfo struct {
	ID   string     `json:"id"`
	Meta []Metainfo `json:"meta"`
}

type Metainfo struct {
	Filename     string `json:"filename"`
	Cmeta        string `json:"custom_meta"`
	Size         int    `json:"size"`
	Content_hash string `json:"content_hash"`
	MetaCustom   *Custom_meta
}

type Custom_meta struct {
	Part_num int `json:"part_num"`
}

func getMeta(body []byte) (*Maininfo, error) {
	var s = new(Maininfo)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func main() {
	res, err := http.Get("http://localhost:5050/v1/file/meta/36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80?path=/&filename=big.txt")
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getMeta([]byte(body))
	fmt.Println("Main stuff", s)

	var metaInfo []Metainfo
	metaInfo = s.Meta
	fmt.Println("metaInfo", metaInfo)
	for i := range metaInfo {
		var cm = new(Custom_meta)
		part := metaInfo[i].Cmeta
		json.Unmarshal([]byte(part), &cm)
		fmt.Println("cm part", cm.Part_num)
		fmt.Println("part", part)
		metaInfo[i].MetaCustom = cm

		if err != nil {
			fmt.Println("err")
		}
		hash := metaInfo[i].Content_hash
		fmt.Println("hash", hash)
	}

	sort.Slice(metaInfo, func(i, j int) bool {
		return metaInfo[i].MetaCustom.Part_num < metaInfo[j].MetaCustom.Part_num
	})

	fmt.Println("entire meta", metaInfo)

	for j := range metaInfo {
		file := metaInfo[j].Filename
		hash := metaInfo[j].Content_hash
		targetUrl := "http://localhost:5050/v1/file/download/36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"
		downloadFile(file, "/", targetUrl, hash)
	}
}

func downloadFile(filename, filepath, targetUrl, part_hash string) (string, error) {
	url := targetUrl + "?path=" + "/" + "&filename=" + filename + "&part_hash=" + part_hash
	fmt.Println("complete url ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	// defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "false", err
	}

	file, _ := os.Create("file" + part_hash + ".txt")
	fmt.Println("after file ")
	io.Copy(file, resp.Body)

	defer file.Close()
	return "true", err

}

// var Messages = make(chan *os.File)
func uploadFile(filename string, reader io.Reader, wg *sync.WaitGroup, meta string) error {
	defer wg.Done()
	bodyReader, bodyWriter := io.Pipe()
	multiWriter := multipart.NewWriter(bodyWriter)
	go func() {
		// fmt.Println("body buffer", bodyWriter)

		// this step is very important

		fileWriter, err := multiWriter.CreateFormFile("uploadFile", filename)
		if err != nil {
			bodyWriter.CloseWithError(err)
			return
		}

		//iocopy
		_, err = io.Copy(fileWriter, reader)
		if err != nil {
			bodyWriter.CloseWithError(err)
			return
		}

		// Create a form field writer for field label
		metaWriter, err := multiWriter.CreateFormField("custom_meta")
		if err != nil {
			bodyWriter.CloseWithError(err)
			return
		}
		metaWriter.Write([]byte(meta))

		bodyWriter.CloseWithError(multiWriter.Close())
	}()
	contentType := multiWriter.FormDataContentType()
	targetUrl := "http://localhost:5050/v1/file/upload/36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"
	resp, err := http.Post(targetUrl, contentType, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err

	}
	fmt.Println("resp", string(resp_body))
	return nil
}

func storeInFile(in io.Reader, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	destfilename := fmt.Sprintf("%s.%d", "big.txt", i)
	fmt.Println("file to be created", destfilename)
	f, err := os.Create(destfilename)
	defer f.Close()
	checkErr(err)
	// copy from reader data into writer file
	_, err = io.Copy(f, in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("file created", destfilename)
}

//ChunkingFilebyShards is used to divide the file in chunks using erasure coding
func (fi *FileInfo) ChunkingFilebyShards() {
	runtime.GOMAXPROCS(30)
	if fi.DataShards > 257 {
		fmt.Fprintf(os.Stderr, "Error: Too many data shards\n")
		os.Exit(1)
	}
	fname := fi.File

	// Create encoding matrix.
	enc, err := reedsolomon.NewStreamC(fi.DataShards, fi.ParShards, true, true)
	checkErr(err)

	fmt.Println("Opening", fname)
	f, err := os.Open(fname)
	checkErr(err)

	instat, err := f.Stat()
	checkErr(err)

	shards := fi.DataShards
	var wg sync.WaitGroup
	wg.Add(18)

	out1 := make([]io.Writer, shards)
	out2 := make([]io.Writer, shards)
	out := make([]io.Writer, shards)

	in := make([]io.Reader, shards)
	inr := make([]io.Reader, shards)
	// Create the resulting files.

	for i := range out {
		outfn := fmt.Sprintf("Part : %d", i)
		meta := fmt.Sprintf("{\"part_num\" : %d}", i)
		fmt.Println("Creating", outfn)
		pr, pw := io.Pipe()
		npr, npw := io.Pipe()
		out1[i] = pw
		out2[i] = npw
		out[i] = io.MultiWriter(pw, npw)
		//out[i] = pw
		checkErr(err)
		//tr := io.TeeReader(pr, f)
		in[i] = pr
		inr[i] = npr
		//destfilename := fmt.Sprintf("%s.%d", "big.txt", i)
		go uploadFile("big.txt", npr, &wg, meta)
		//go storeInFile(npr, i, &wg);
	}

	// Create parity output writers
	parity := make([]io.Writer, 6)
	for i := range parity {
		// destfilename := fmt.Sprintf("%s.%d", "big.txt", 10+i)
		// fmt.Println("file to be created" , destfilename)
		// f, err := os.Create(destfilename)
		// defer f.Close()
		// checkErr(err)
		// parity[i] = f
		// //parity[i] = out[10+i]
		// //defer out[10+i].(*io.PipeWriter).Close()
		// fmt.Println("file created" , destfilename)
		pr, pw := io.Pipe()
		parity[i] = pw
		//destfilename := fmt.Sprintf("%s.%d", "big.txt", 10+i)
		meta := fmt.Sprintf("{\"part_num\" : %d}", 10+i)
		go uploadFile("big.txt", pr, &wg, meta)

	}

	go func() {
		defer wg.Done()
		// Encode parity
		err = enc.Encode(in, parity)
		checkErr(err)
		for i := range parity {
			parity[i].(*io.PipeWriter).Close()
		}
	}()

	go func() {
		defer wg.Done()
		// Do the split
		err = enc.Split(f, out, instat.Size())
		checkErr(err)
		fmt.Println("Done with split")
		for i := range out {
			out2[i].(*io.PipeWriter).Close()
			out1[i].(*io.PipeWriter).Close()
			//out2[i].(*io.PipeWriter).Close()
			//out[i].(*io.PipeWriter).Close()

		}
	}()

	wg.Wait()

	fmt.Printf("File split into %d data + %d parity shards.\n", 10, 6)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}

func (fi *FileInfo) DecodeFileByShards() {

	fname := fi.File

	// Create matrix
	enc, err := reedsolomon.NewStream(fi.DataShards, fi.ParShards)
	checkErr(err)

	// Open the inputs
	shards, size, err := openInput(fi.DataShards, fi.ParShards, fname)
	checkErr(err)

	// Verify the shards
	ok, err := enc.Verify(shards)
	// fmt.Println("verify", ok)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		shards, size, err = openInput(fi.DataShards, fi.ParShards, fname)
		checkErr(err)
		// Create out destination writers
		out := make([]io.Writer, len(shards))
		for i := range out {

			if shards[i] == nil {
				outfn := fmt.Sprintf("%s", fname)
				fmt.Println("Creating", outfn)

				out[i], err = os.Create(outfn)
				checkErr(err)

			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		// Close output.
		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				checkErr(err)
			}
		}

		shards, size, err = openInput(fi.DataShards, fi.ParShards, fname)
		ok, err = enc.Verify(shards)
		fmt.Println("ok", ok)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted:", err)
			os.Exit(1)
		}
		checkErr(err)
	}

	// Join the shards and write them
	outfn := fi.OutDir
	if outfn == "" {
		outfn = fname
	}

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	checkErr(err)

	shards, size, err = openInput(fi.DataShards, fi.ParShards, fname)
	checkErr(err)

	err = enc.Join(f, shards, int64(fi.DataShards)*size)
	checkErr(err)
}

func openInput(dataShards, parShards int, fname string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		ofn := fmt.Sprintf("%d", i+1)
		// targetUrl := "http://localhost:5050/v1/file/meta/36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"
		// // downloadFile(infn, "/", targetUrl, ofn)

		// filename, _ := downloadFile(fname, "/", targetUrl)

		f, err := os.Open(ofn)
		if err != nil {
			fmt.Println("Error reading filesssssss", err)
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		checkErr(err)
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}

	}

	return shards, size, nil

}
