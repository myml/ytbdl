// ytbdl project main.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-macaron/macaron"
	"github.com/otium/ytdl"
)

type Video struct {
	Title     string
	Thumbnail string
	Formats   []Format
}
type Format struct {
	Res  string
	Ext  string
	Url  string
	Clen string
}

func main() {
	//	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:7070", nil, proxy.Direct)
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	http.DefaultClient.Transport = &http.Transport{Dial: dialer.Dial}

	m := macaron.Classic()
	m.Use(macaron.Static("html"))
	m.Use(macaron.Renderer(macaron.RenderOptions{IndentJSON: true}))
	m.Get("/video/:id", func(ctx *macaron.Context) {
		id := ctx.Params("id")
		log.Println(id)
		info, err := ytdl.GetVideoInfoFromID(id)
		if err != nil {
			log.Panic(err)
		}
		formats := info.Formats.
			Filter(ytdl.FormatExtensionKey, []interface{}{"mp4"}).
			Filter(ytdl.FormatAudioEncodingKey, []interface{}{""})
		out := Video{
			Title:     info.Title,
			Thumbnail: info.GetThumbnailURL(ytdl.ThumbnailQualityDefault).String(),
			Formats:   make([]Format, len(formats)),
		}
		for i := range formats {
			out.Formats[i] = Format{
				Res: formats[i].Resolution,
				Ext: formats[i].Extension,
				Url: func() string {
					u, err := info.GetDownloadURL(formats[i])
					if err != nil {
						log.Panic(err)
					}
					return u.String()
				}(),
				Clen: formats[i].ValueForKey("clen").(string),
			}
		}
		ctx.JSON(200, out)
	})
	m.Get("dl/:fname", func(ctx *macaron.Context) {
		fname := ctx.Params("fname")
		u := ctx.Query("url")
		clen := ctx.Query("clen")
		log.Println(u, clen)
		resp, err := http.Get(u)
		if err != nil {
			log.Panic(err)
		}
		ctx.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
		ctx.Header().Set("Content-type", "application/octet-stream")
		ctx.Header().Set("Content-Length", clen)
		log.Println(io.Copy(ctx, resp.Body))
	})
	m.Run()
}