package main

/*

	GifCreator - Telegram Bot to create GIFs from a series of images.
	Copyright (C) 2021  MassiveBox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"

	"image/gif"
	"image/jpeg"
	"net/http"
	"os"
)

type savedFrame struct {
	FileID string `json:"file_id"`
	Length int    `json:"length"`
}

func savePhoto(fileID string, duration int, userID int64, db *sql.DB) error {

	data, err := databaseGetUserData(db, userID)
	if err != nil {
		return err
	}

	var frames []savedFrame
	if data.Frames != "" {
		err = json.Unmarshal([]byte(data.Frames), &frames)
		if err != nil {
			return err
		}
	}

	if len(frames) > 150 {
		return errors.New("limits")
	}

	if !data.HasVoted {
		var totalFrameCount int
		for _, frame := range frames {
			totalFrameCount += frame.Length
		}
		if totalFrameCount > 1000 {
			return errors.New("free")
		}
	}

	frames = append(frames, savedFrame{
		FileID: fileID,
		Length: duration,
	})

	marshaledFrames, err := json.Marshal(frames)
	if err != nil {
		return err
	}

	data.Frames = string(marshaledFrames)
	databaseUpdateUserData(db, data)

	return nil

}

func (pt passthrough) downloadFileID(ID string) (image.Image, error) {

	fileURL, err := pt.Bot.GetFileDirectURL(ID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return jpeg.Decode(resp.Body)

}

func (pt passthrough) renderGIF(data userData) error {

	var frames []savedFrame
	if data.Frames != "" {
		err := json.Unmarshal([]byte(data.Frames), &frames)
		if err != nil {
			return err
		}
	}else{
		return errors.New("no_frames")
	}

	var (
		imagesProv []image.Image
		maxW, maxH int
	)

	for _, frame := range frames {

		img, err := pt.downloadFileID(frame.FileID)
		if err != nil {
			return err
		}
		if img.Bounds().Max.X > maxW {
			maxW = img.Bounds().Max.X
		}
		if img.Bounds().Max.Y > maxH {
			maxH = img.Bounds().Max.Y
		}
		imagesProv = append(imagesProv, img)

	}

	outGif := &gif.GIF{}

	for key, imageProv := range imagesProv {

		differenceW := maxW - imageProv.Bounds().Max.X
		differenceH := maxH - imageProv.Bounds().Max.Y

		var pt0x, pt0y int
		pt0x -= differenceW / 2
		pt0y -= differenceH / 2

		palettedImage := image.NewPaletted(image.Rect(0,0,maxW,maxH), palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, imageProv, image.Pt(pt0x, pt0y), draw.Over)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, frames[key].Length)

	}

	out, err := os.Create("tmp/"+fmt.Sprintf("%d.gif", data.UserID))
	if err != nil {
		return err
	}
	defer out.Close()

	err = gif.EncodeAll(out, outGif)
	if err != nil {
		return err
	}

	return nil

}