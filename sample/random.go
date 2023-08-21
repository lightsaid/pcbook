package sample

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lightsaid/pcbook/pb"
)

var myRand *rand.Rand

func init() {
	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// randomStringFromSet 随机返回其中一个string
func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}

	return a[myRand.Intn(n)]
}

// randomBool 返回 true/false
func randomBool() bool {
	return myRand.Intn(2) == 1
}

// randomInt 随机返回一个整数，介于 [min, max]
func randomInt(min, max int) int {
	return min + myRand.Int()%(max-min+1)
}

// randomFloat64 返回一个float64浮点数，介于[min, max]
func randomFloat64(min, max float64) float64 {
	return min + myRand.Float64()*(max-min)
}

// randomFloat32 返回一个float32浮点数，介于[min, max]
func randomFloat32(min, max float32) float32 {
	return min + myRand.Float32()*(max-min)
}

// randomID 随机返回uuid
func randomID() string {
	return uuid.New().String()
}

// randomKeyboardLayout 随机返回键盘布局
func randomKeyboardLayout() pb.Keyboard_Layout {
	switch myRand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

// randomScreenPanel 随机返回显示器材质
func randomScreenPanel() pb.Screen_Panel {
	if myRand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

// randomCPUBrand 从预设值中随机返回 CPU 品牌
func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

// randomCPUName 根据品牌名从预设值中随机返回名称
func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon E-2286M",
			"Core i9-9980HK",
			"Core i7-9750H",
			"Core i5-9400F",
			"Core i3-1005G1",
		)
	}

	return randomStringFromSet(
		"Ryzen 7 PRO 2700U",
		"Ryzen 5 PRO 3500U",
		"Ryzen 3 PRO 3200GE",
	)
}

// randomGPUBrand 随机返回 GPU 品牌
func randomGPUBrand() string {
	return randomStringFromSet("Nvidia", "AMD")
}

// randomGPUName 根据品牌名从预设值中随机返回名称
func randomGPUName(brand string) string {
	if brand == "Nvidia" {
		return randomStringFromSet(
			"RTX 2060",
			"RTX 2070",
			"GTX 1660-Ti",
			"GTX 1070",
		)
	}

	return randomStringFromSet(
		"RX 590",
		"RX 580",
		"RX 5700-XT",
		"RX Vega-56",
	)
}

// randomLaptopBrand 从预设值中随机返回笔记本品牌
func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}

// randomLaptopName 从预设值中随机返回笔记本名称
func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1", "Thinkpad P53")
	}
}
