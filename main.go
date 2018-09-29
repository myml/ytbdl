// ytbdl project main.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/proxy"

	"github.com/go-macaron/macaron"
	"github.com/otium/ytdl"
)

type Video struct {
	ID      string
	Title   string
	Formats []Format
}
type Format struct {
	Itag int
	Res  string
	Ext  string
	Clen string
}

func main() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:7070", nil, proxy.Direct)
	if err != nil {
		log.Panic(err)
	}
	http.DefaultClient.Transport = &http.Transport{Dial: dialer.Dial}

	m := macaron.Classic()
	m.Use(macaron.Static("html"))
	m.Use(macaron.Renderer(macaron.RenderOptions{IndentJSON: true}))
	m.Get("/video/:id.jpg", func(ctx *macaron.Context) {
		id := ctx.Params("id")
		uri := fmt.Sprintf("http://img.youtube.com/vi/%s/default.jpg", id)
		resp, err := http.Get(uri)
		if err != nil {
			log.Panic(err)
		}
		io.Copy(ctx, resp.Body)
	})
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
			ID:      id,
			Title:   info.Title,
			Formats: make([]Format, len(formats)),
		}
		for i := range formats {
			out.Formats[i] = Format{
				Itag: formats[i].Itag,
				Res:  formats[i].Resolution,
				Ext:  formats[i].Extension,
				Clen: formats[i].ValueForKey("clen").(string),
			}
		}
		ctx.JSON(200, out)
	})
	m.Get("/video/:id/format/:itag", func(ctx *macaron.Context) {
		id := ctx.Params("id")
		itag := ctx.ParamsInt("itag")
		log.Println(id, itag)
		info, err := ytdl.GetVideoInfoFromID(id)
		if err != nil {
			log.Panic(err)
		}
		for i := range info.Formats {
			if info.Formats[i].Itag == itag {
				ctx.Header().Set("Content-type", "application/octet-stream")
				fname := fmt.Sprintf("%s_%s.%s", info.Title,
					info.Formats[i].Resolution,
					info.Formats[i].Extension)
				ctx.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fname))

				uri, err := info.GetDownloadURL(info.Formats[i])
				if err != nil {
					log.Panic(err)
				}
				resp, err := http.Get(uri.String())
				if err != nil {
					log.Panic(err)
				}
				defer resp.Body.Close()
				ctx.Header().Set("Content-Length", fmt.Sprint(resp.ContentLength))
				log.Println(io.Copy(ctx.Resp, resp.Body))
			}
		}
	})
	m.Run(4444)
}
