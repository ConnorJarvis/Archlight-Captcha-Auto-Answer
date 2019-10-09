package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
)

var textColor color.RGBA

func main() {
	bounds := screenshot.GetDisplayBounds(0)
	answerIndex := 0
	for {

		img, err := screenshot.CaptureRect(image.Rect(bounds.Dx()/3, bounds.Dy()/3, (bounds.Dx()/3)*2, (bounds.Dy()/3)*2))
		if err != nil {
			fmt.Println(err)
			goto End
		}
		answerIndex, err = detectQuestion(img)
		if err != nil {
			fmt.Println(err)
			goto End
		}
		if err == nil {
			for i := 0; i < answerIndex; i++ {
				robotgo.KeyTap("down")
				time.Sleep(time.Millisecond * 10)
			}
			time.Sleep(time.Millisecond * 20)
			robotgo.KeyTap("enter")
			time.Sleep(time.Millisecond * 2000)
		}
	End:
		time.Sleep(time.Millisecond * 200)

	}

}

func checkPoint(x, y int, img image.Image) bool {
	// for xCheck := x; xCheck < x+2; xCheck++ {
	// 	for yCheck := y; yCheck < y+2; yCheck++ {
	r, g, b, _ := img.At(x, y).RGBA()
	rV, gV, bV := 65535, 65535, 65535
	margin := 19000
	if !(((int(rV)-margin) < int(r) && int(r) < (int(rV)+margin)) && ((int(gV)-margin) < int(g) && int(g) < (int(gV)+margin)) && ((int(bV)-margin) < int(b) && int(b) < (int(bV)+margin))) {
		return false
	}
	// if r != 65535 || g != 65535 || b != 65535 {
	// 	r
	// }
	// 	}
	// }
	return true
}

func detectAnswers(src image.Image) string {

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	text := ""
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if checkPoint(x, y, src) {
				if detectAnswer0(x, y, src) {
					text += "0"
					x += 4
				} else if detectAnswer8(x, y, src) {
					text += "8"
					x += 4
				} else if detectNumber2(x, y, src) {
					text += "2"
					x += 4
				} else if detectNumber6(x, y, src) {
					text += "6"
					x += 4
				} else if detectNumber4(x, y, src) {
					text += "4"
					x += 4
				} else if detectNumber1(x, y, src) {
					text += "1"
					x += 4
				} else if detectNumber5(x, y, src) {
					text += "5"
					x += 4
				} else if detectNumber3(x, y, src) {
					text += "3"
					x += 4
				} else if detectNumber9(x, y, src) {
					text += "9"
					x += 4
				} else if detectNumber7(x, y, src) {
					text += "7"
					x += 4
				} else if detectPlus(x, y, src) {
					text += "+"
					x += 4
				} else if detectDash(x, y, src) {
					text += "-"
					x += 4
				}
			}
		}
	}

	return text
}

func detectQuestion(img image.Image) (int, error) {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if checkPoint(x, y, img) {
				if detectPlease(x, y, img) {

					rgbImg := img.(*image.RGBA)
					subImg := rgbImg.SubImage(image.Rect(x, y+28, x+60, y+8))
					equation := readEquation(subImg)

					if equation != "" {
						expression, err := govaluate.NewEvaluableExpression(equation)
						if err != nil {
							fmt.Println(err)
							return 0, errors.New("No answer")
						}
						result, err := expression.Evaluate(nil)
						if err != nil {
							fmt.Println(err)
							return 0, errors.New("No answer")
						}

						answers := []int{}
						for i := 0; i < 5; i++ {

							answersImg := rgbImg.SubImage(image.Rect(x, y+29+(i*14), x+17, y+43+(i*14)))
							answer, err := strconv.Atoi(detectAnswers(answersImg))
							if err != nil {
								fmt.Println(err)
								return 0, errors.New("No answer")
							}
							answers = append(answers, answer)
						}
						for i := 0; i < len(answers); i++ {
							if answers[i] == int(result.(float64)) {
								return i, nil
							}
						}
					}

					return 0, errors.New("No question")
				}
			}
		}
	}
	// for yCheck := y + 1; yCheck <= y+5; yCheck += 1 {
	// 	if !checkPoint(x, yCheck, img) {
	// 		return false
	// 	}
	// }
	// for yCheck := y; yCheck <= y+5; yCheck += 1 {
	// 	if !checkPoint(x+4, yCheck, img) {
	// 		return false
	// 	}
	// }
	// for xCheck := x + 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y-1, img) {

	// 		return false
	// 	}
	// }

	// for xCheck := x + 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+6, img) {
	// 		return false
	// 	}
	// }
	// return true
	return 0, errors.New("No question")
}

func detectPlease(x, y int, img image.Image) bool {
	for xCheck := x; xCheck <= x+1; xCheck += 1 {
		for yCheck := y; yCheck <= y+7; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	for xCheck := x + 7; xCheck <= x+8; xCheck += 1 {
		for yCheck := y - 1; yCheck <= y+7; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	for xCheck := x + 11; xCheck <= x+15; xCheck += 1 {

		if !checkPoint(xCheck, y+4, img) {
			return false
		}

	}
	for xCheck := x - 2; xCheck <= x-1; xCheck += 1 {
		for yCheck := y; yCheck <= y+7; yCheck += 1 {
			if checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	return true
}

func readEquation(src image.Image) string {
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	text := ""
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if checkPoint(x, y, src) {
				if detectNumber2(x, y, src) {
					text += "2"
					x += 8
				} else if detectNumber6(x, y, src) {
					text += "6"
					x += 8
				} else if detectNumber4(x, y, src) {
					text += "4"
					x += 8
				} else if detectNumber1(x, y, src) {
					text += "1"
					x += 8
				} else if detectNumber5(x, y, src) {
					text += "5"
					x += 8
				} else if detectNumber3(x, y, src) {
					text += "3"
					x += 8
				} else if detectNumber9(x, y, src) {
					text += "9"
					x += 8
				} else if detectNumber7(x, y, src) {
					text += "7"
					x += 8
				} else if detectPlus(x, y, src) {
					text += "+"
					x += 8
				} else if detectDash(x, y, src) {
					text += "-"
					x += 8
				} else if detectNumber8(x, y, src) {
					text += "8"
					x += 8
				}
			}
		}
	}

	return text
}

func detectAnswer8(x, y int, img image.Image) bool {
	for xCheck := x; xCheck <= x+1; xCheck += 1 {
		for yCheck := y; yCheck <= y+1; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}

	for xCheck := x; xCheck <= x+1; xCheck += 1 {
		for yCheck := y + 3; yCheck <= y+5; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}

	for xCheck := x + 4; xCheck <= x+1; xCheck += 1 {
		for yCheck := y; yCheck <= y+1; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}

	for xCheck := x + 4; xCheck <= x+1; xCheck += 1 {
		for yCheck := y + 3; yCheck <= y+5; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	// for xCheck := x - 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+3, img) {
	// 		return false
	// 	}
	// }
	return true
}

func detectAnswer0(x, y int, img image.Image) bool {
	for xCheck := x; xCheck <= x+1; xCheck += 1 {
		for yCheck := y; yCheck <= y+5; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}

	for xCheck := x + 4; xCheck <= x+1; xCheck += 1 {
		for yCheck := y; yCheck <= y+5; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	// for xCheck := x - 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+3, img) {
	// 		return false
	// 	}
	// }
	return true
}

func detectNumber1(x, y int, img image.Image) bool {

	for xCheck := x + 1; xCheck <= x+2; xCheck += 1 {
		for yCheck := y; yCheck <= y+6; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	for xCheck := x; xCheck <= x+3; xCheck += 1 {

		if !checkPoint(xCheck, y+6, img) {
			return false
		}

	}
	if !checkPoint(x+2, y-1, img) {
		return false
	}
	// for xCheck := x - 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+3, img) {
	// 		return false
	// 	}
	// }
	return true
}

func detectNumber6(x, y int, img image.Image) bool {

	for yCheck := y; yCheck <= y+5; yCheck += 1 {
		if !checkPoint(x, yCheck, img) {
			return false
		}
	}

	for yCheck := y + 3; yCheck <= y+5; yCheck += 1 {
		if !checkPoint(x+3, yCheck, img) {
			return false
		}
	}

	for xCheck := x + 2; xCheck <= x+3; xCheck += 1 {

		if !checkPoint(xCheck, y-1, img) {
			return false
		}

	}

	// for xCheck := x - 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+3, img) {
	// 		return false
	// 	}
	// }
	return true
}

func detectNumber2(x, y int, img image.Image) bool {

	for yCheck := y; yCheck <= y+2; yCheck += 1 {
		if !checkPoint(x+4, yCheck, img) {
			return false
		}
	}

	for xCheck := x; xCheck <= x+5; xCheck += 1 {

		if !checkPoint(xCheck, y+6, img) {
			return false
		}

	}
	if !checkPoint(x+3, y+3, img) {
		return false
	}
	if !checkPoint(x+2, y+4, img) {
		return false
	}
	// for xCheck := x - 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+3, img) {
	// 		return false
	// 	}
	// }
	return true
}

func detectNumber4(x, y int, img image.Image) bool {

	for xCheck := x; xCheck <= x+5; xCheck += 1 {

		if !checkPoint(xCheck, y, img) {
			return false
		}

	}
	for xCheck := x + 3; xCheck <= x+4; xCheck += 1 {
		for yCheck := y - 3; yCheck <= y+2; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	// for xCheck := x - 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+3, img) {
	// 		return false
	// 	}
	// }
	return true
}

func detectNumber8(x, y int, img image.Image) bool {

	for yCheck := y + 1; yCheck <= y+6; yCheck += 1 {
		if !checkPoint(x-1, yCheck, img) {
			return false
		}
	}
	for yCheck := y + 1; yCheck <= y+2; yCheck += 1 {
		if !checkPoint(x+2, yCheck, img) {
			return false
		}
	}
	for yCheck := y + 4; yCheck <= y+6; yCheck += 1 {
		if !checkPoint(x+2, yCheck, img) {
			return false
		}
	}
	// for xCheck := x - 1; xCheck <= x+3; xCheck += 1 {
	// 	if !checkPoint(xCheck, y+3, img) {
	// 		return false
	// 	}
	// }
	return true
}

func detectNumber5(x, y int, img image.Image) bool {

	for xCheck := x; xCheck <= x+1; xCheck += 1 {
		for yCheck := y; yCheck <= y+3; yCheck += 1 {
			if !checkPoint(xCheck, yCheck, img) {
				return false
			}
		}
	}
	for xCheck := x + 2; xCheck <= x+4; xCheck += 1 {

		if !checkPoint(xCheck, y, img) {
			return false
		}

	}
	for yCheck := y + 4; yCheck <= y+6; yCheck += 1 {
		if !checkPoint(x+3, yCheck, img) {
			return false
		}
	}

	return true
}

func detectNumber7(x, y int, img image.Image) bool {

	for xCheck := x; xCheck <= x+5; xCheck += 1 {
		if !checkPoint(xCheck, y, img) {
			return false
		}
	}

	for yCheck := y + 1; yCheck <= y+2; yCheck += 1 {
		if !checkPoint(x+4, yCheck, img) {
			return false
		}
	}
	for yCheck := y + 3; yCheck <= y+4; yCheck += 1 {
		if !checkPoint(x+3, yCheck, img) {
			return false
		}
	}
	for yCheck := y + 5; yCheck <= y+6; yCheck += 1 {
		if !checkPoint(x+2, yCheck, img) {
			return false
		}
	}
	if !checkPoint(x+1, y+7, img) {
		return false
	}
	return true
}

func detectNumber3(x, y int, img image.Image) bool {

	for yCheck := y + 1; yCheck <= y+2; yCheck += 1 {
		if !checkPoint(x+2, yCheck, img) {
			return false
		}
	}
	for yCheck := y + 4; yCheck <= y+6; yCheck += 1 {
		if !checkPoint(x+2, yCheck, img) {
			return false
		}
	}
	for xCheck := x; xCheck <= x+1; xCheck += 1 {
		if !checkPoint(xCheck, y+3, img) {
			return false
		}
	}
	if !checkPoint(x-2, y+1, img) {
		return false
	}
	return true
}

func detectNumber9(x, y int, img image.Image) bool {

	for yCheck := y; yCheck <= y+2; yCheck += 1 {
		if !checkPoint(x, yCheck, img) {
			return false
		}
	}
	if !checkPoint(x+1, y+3, img) {
		return false
	}

	for yCheck := y; yCheck <= y+5; yCheck += 1 {
		if !checkPoint(x+3, yCheck, img) {
			return false
		}
	}
	if !checkPoint(x+1, y+6, img) {
		return false
	}

	return true
}

func detectDash(x, y int, img image.Image) bool {

	for xCheck := x; xCheck <= x+3; xCheck += 1 {
		if !checkPoint(xCheck, y, img) {
			return false
		}
	}
	for xCheck := x; xCheck <= x+3; xCheck += 1 {
		if checkPoint(xCheck, y-1, img) {
			return false
		}
	}
	for xCheck := x; xCheck <= x+3; xCheck += 1 {
		if checkPoint(xCheck, y+1, img) {
			return false
		}
	}
	return true
}

func detectPlus(x, y int, img image.Image) bool {

	for xCheck := x; xCheck <= x+6; xCheck += 1 {
		if !checkPoint(xCheck, y, img) {
			return false
		}
	}
	for yCheck := y - 3; yCheck <= y+3; yCheck += 1 {
		if !checkPoint(x+3, yCheck, img) {
			return false
		}
	}
	return true
}
