package resize

import (
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

type Progress func(progress, total int)

type Options struct {
	Quality   int
	ImageSize int
	ThumbSize int
}

func Resize(dir string, opts *Options, prog Progress) error {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return e
	}

	p := 0
	t := len(files)
	c := make(chan bool)
	defer close(c)

	for _, f := range files {
		go process(dir, f.Name(), opts, c)
	}

	for p < t {
		<-c
		p++
		prog(p, t)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Millisecond * 200)
	prog(t, t)

	return OpenDir(filepath.Join(dir, "output"))
}

func process(dir, name string, opts *Options, c chan bool) {
	defer func() { c <- true }()

	fileIn, e := os.Open(filepath.Join(dir, name))
	if e != nil {
		log.Println("err - open input file:", e)
		return
	}
	defer fileIn.Close()

	fileOut, e := mkFile(dir, "output", name)
	if e != nil {
		log.Println("err - create output file:", e)
		return
	}
	defer fileOut.Close()

	thumbOut, e := mkFile(dir, "output", "thumbnail", name)
	if e != nil {
		log.Println("err - create thumbnail file:", e)
		return
	}
	defer thumbOut.Close()

	e = resize(fileIn, fileOut, thumbOut, opts)
	if e != nil {
		log.Println("err - resize files:", e)
		return
	}
}

func resize(in io.Reader, imgOut, thumbOut io.Writer, opts *Options) error {
	i, _, e := image.Decode(in)
	if e != nil {
		return e
	}

	img := imaging.Fit(i, opts.ImageSize, opts.ImageSize, imaging.Lanczos)
	thumb := imaging.Fit(i, opts.ThumbSize, opts.ThumbSize, imaging.Lanczos)

	e = jpeg.Encode(imgOut, img, &jpeg.Options{Quality: opts.Quality})
	if e != nil {
		return e
	}

	return jpeg.Encode(thumbOut, thumb, &jpeg.Options{Quality: opts.Quality})
}

func mkFile(path ...string) (*os.File, error) {
	file := filepath.Join(path...)
	e := os.MkdirAll(filepath.Dir(file), os.ModePerm)
	if e != nil {
		return nil, e
	}
	return os.Create(file)
}
