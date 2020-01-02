// Copyright 2019-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fb

import (
  "image"
  "io/ioutil"
)

const width = 640
const height = 480
const bpp = 3

func DrawImageAt(img image.Image, posx int, posy int) {
  var buf []byte
  buf = make([]byte, width*height*bpp)
  DrawOnBufAt(buf, img, posx, posy, width, bpp)
  ioutil.WriteFile("/dev/fb0", buf, 0600)
}

func DrawOnBufAt(
  buf []byte,
  img image.Image,
  posx int,
  posy int,
  width int,
  bpp int,
) {
  for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
    for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
      r, g, b, a := img.At(x, y).RGBA()
      offset := bpp * ((posy+y)*width + posx+x)
      // framebuffer is BGR(A)
      buf[offset + 0] = byte(b)
      buf[offset + 1] = byte(g)
      buf[offset + 2] = byte(r)
      if (bpp == 4) {
        buf[offset + 3] = byte(a)
      }
    }
  }
}

func DrawScaledOnBufAt(
  buf []byte,
  img image.Image,
  posx int,
  posy int,
  factor int,
  width int,
  bpp int,
) {
  for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
    for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
      r, g, b, a := img.At(x, y).RGBA()
      for sx := 1; sx <= factor; sx++ {
        for sy := 1; sy <= factor; sy++ {
          offset := (posx+x*factor + (posy+y*factor+sy)*width + sx) * bpp
          buf[offset + 0] = byte(b)
          buf[offset + 1] = byte(g)
          buf[offset + 2] = byte(r)
          if (bpp == 4) {
            buf[offset + 3] = byte(a)
          }
        }
      }
    }
  }
}

func DrawScaledImageAt(img image.Image, posx int, posy int, factor int) {
  var buf []byte
  buf = make([]byte, width*height*bpp)
  DrawScaledOnBufAt(buf, img, posx, posy, factor, width, bpp)
  size := 3
  for digit := 0; digit < 9; digit++ {
    DrawDigitAt(buf, digit, 205 + digit*15, 130, width, bpp, size)
  }
  ioutil.WriteFile("/dev/fb0", buf, 0600)
}
