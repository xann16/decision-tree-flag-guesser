package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"../pkg/flagger"
)

func main() {
	imgPath := "../static/img/flags/"
	dstPath := "../temp/kb_traits_raw.txt"
	if len(os.Args) > 1 {
		imgPath = os.Args[1]
	}
	tags := loadTags(imgPath)
	data := initTraitMap(tags)

	gatherTraits(tags, data)

	file, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data.WriteOrdered(file, true)
}

func loadTags(basePath string) []string {
	var tags []string
	err := filepath.Walk(basePath,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".png" {
				tag := filepath.Base(path)
				tag = tag[:len(tag)-4]
				tags = append(tags, tag)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
	return tags
}

func initTraitMap(tags []string) flagger.TraitMap {
	return flagger.CreateTraitMap(tags)
}

func gatherTraits(tags []string, tm flagger.TraitMap) {
	total := len(tags)
	for i, tag := range tags {
		fmt.Printf("[%3d/%3d] Processing tag %v:\n", i+1, total, tag)
		var traits uint64

		name := getName()
		traits += processColours()
		traits += processLayout()
		traits += processSymmetry()
		traits += processFeatures()

		printCheck(tag, name, traits)
		tm.Update(tag, name, traits)
	}
}

func printCheck(tag, name string, traits uint64) {
	fmt.Printf("\n(%s) %s:\n", tag, name)
	for i := uint64(0); i < 64; i++ {
		if (traits & (1 << i)) != 0 {
			fmt.Printf(" -> %s\n", RULES[i])
		}
	}
}

func getName() string {
	var line string
	fmt.Printf("name: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	return strings.TrimSpace(line)
}

func processColours() (flags uint64) {
	var line string
	fmt.Printf("colours: ")
	fmt.Scanln(&line)

	if strings.Contains(line, "w") {
		flags += uint64(0x01)
	}
	if strings.Contains(line, "b") {
		flags += uint64(0x02)
	}
	if strings.Contains(line, "y") {
		flags += uint64(0x04)
	}
	if strings.Contains(line, "o") {
		flags += uint64(0x08)
	}
	if strings.Contains(line, "r") {
		flags += uint64(0x10)
	}
	if strings.Contains(line, "u") {
		flags += uint64(0x20)
	}
	if strings.Contains(line, "l") {
		flags += uint64(0x40)
	}
	if strings.Contains(line, "g") {
		flags += uint64(0x80)
	}

	return
}

func processLayout() (flags uint64) {
	var line string
	isH, isV, isD := false, false, false

	fmt.Printf("stripes: ")
	fmt.Scanln(&line)

	if strings.Contains(line, "h") {
		isH = true
	}
	if strings.Contains(line, "v") {
		isV = true
	}
	if strings.Contains(line, "d") {
		isD = true
	}

	if isH {
		flags += uint64(0x0100)
	}
	if isV {
		flags += uint64(0x0200)
	}
	if isH && isV {
		flags += uint64(0x0400)
	}
	if isD {
		flags += uint64(0x0800)
	}

	if isH || isV {
		if getYesNo("all_different_colours: ") {
			flags += uint64(0x1000)
		}
		if getYesNo("alternating_two_colours: ") {
			flags += uint64(0x2000)
		}
		if getYesNo("uneven: ") {
			flags += uint64(0x4000)
		}
	}

	if isH {
		count := getCount("horizontal_count: ")
		if count >= 3 {
			flags += uint64(0x008000)
		}
		if count >= 4 {
			flags += uint64(0x010000)
		}
	}

	if isV {
		count := getCount("vertical_count: ")
		if count >= 3 {
			flags += uint64(0x020000)
		}
	}

	if getYesNo("cross: ") {
		flags += uint64(0x040000)

		if getYesNo(" - diagonal: ") {
			flags += uint64(0x080000)
		} else if getYesNo(" - greek: ") {
			flags += uint64(0x100000)
		} else if getYesNo(" - nordic: ") {
			flags += uint64(0x200000)
		}

		if getYesNo(" - border: ") {
			flags += uint64(0x400000)
		}
	}

	if getYesNo("mono_bkg: ") {
		flags += uint64(0x800000)
	}

	if getYesNo("circle: ") {
		flags += uint64(0x01000000)
	}

	return
}

func processSymmetry() (flags uint64) {
	var line string
	isH, isV := false, false

	fmt.Printf("symmetries: ")
	fmt.Scanln(&line)

	if strings.Contains(line, "h") {
		isH = true
	}
	if strings.Contains(line, "v") {
		isV = true
	}

	if isH || isV {
		flags += uint64(0x02000000)
	}
	if isH {
		flags += uint64(0x04000000)
	}
	if isV {
		flags += uint64(0x08000000)
	}
	if isH && isV {
		flags += uint64(0x10000000)
	}

	return
}

func processFeatures() (flags uint64) {
	if getYesNo("features: ") {
		isStar, isMoon := false, false

		if getYesNo(" - stars: ") {
			flags += uint64(0x20000000)
			isStar = true

			if getYesNo("   - white: ") {
				flags += uint64(0x40000000)
			}
			if getYesNo("   - yellow: ") {
				flags += uint64(0x80000000)
			}

			count := getCount("   - count: ")
			if count >= 2 {
				flags += uint64(0x0100000000)
			}
			if count >= 3 {
				flags += uint64(0x0200000000)
			}
			if count >= 4 {
				flags += uint64(0x0400000000)
			}
			if count >= 5 {
				flags += uint64(0x0800000000)
			}
			if count >= 6 {
				flags += uint64(0x1000000000)
			}
			if count >= 10 {
				flags += uint64(0x2000000000)
			}
		}

		if getYesNo(" - sun: ") {
			flags += uint64(0x4000000000)

			if getYesNo("   - half-disk: ") {
				flags += uint64(0x8000000000)
			}
		}

		if getYesNo(" - crescent-moon: ") {
			flags += uint64(0x010000000000)
			isMoon = true

			if getYesNo("   - half-disk: ") {
				flags += uint64(0x020000000000)
			}
		}

		if isStar && isMoon {
			flags += uint64(0x040000000000)
		}

		if getYesNo(" - writings: ") {
			flags += uint64(0x080000000000)

			if getYesNo("   - arabic: ") {
				flags += uint64(0x100000000000)
			}
			if getYesNo("   - nin-latin: ") {
				flags += uint64(0x200000000000)
			}
			if getYesNo("   - nin-english: ") {
				flags += uint64(0x400000000000)
			}
			if getYesNo("   - nin-esp-por: ") {
				flags += uint64(0x800000000000)
			}
		}

		if getYesNo(" - plants: ") {
			flags += uint64(0x01000000000000)

			if getYesNo("   - tree: ") {
				flags += uint64(0x02000000000000)
			}
		}

		if getYesNo(" - weapons: ") {
			flags += uint64(0x04000000000000)

			if getYesNo("   - spear or lance: ") {
				flags += uint64(0x08000000000000)
			}
		}

		if getYesNo(" - tools: ") {
			flags += uint64(0x10000000000000)
		}

		if getYesNo(" - beasts: ") {
			flags += uint64(0x20000000000000)

			if getYesNo("   - mythical: ") {
				flags += uint64(0x40000000000000)
			}
			if getYesNo("   - bird: ") {
				flags += uint64(0x80000000000000)
			}
		}

		if getYesNo(" - uk-flag: ") {
			flags += uint64(0x0100000000000000)
		}
		if getYesNo(" - southern-cross: ") {
			flags += uint64(0x0200000000000000)
		}
		if getYesNo(" - naval: ") {
			flags += uint64(0x0400000000000000)
		}
		if getYesNo(" - coat-of-arms: ") {
			flags += uint64(0x0800000000000000)
		}
		if getYesNo(" - crown: ") {
			flags += uint64(0x1000000000000000)
		}
		if getYesNo(" - human-figures: ") {
			flags += uint64(0x2000000000000000)
		}
		if getYesNo(" - buildings: ") {
			flags += uint64(0x4000000000000000)
		}
		if getYesNo(" - waves: ") {
			flags += uint64(0x8000000000000000)
		}

	}

	return
}

func getYesNo(label string) bool {
	var line string
	fmt.Printf("%s", label)
	fmt.Scanln(&line)
	if strings.Contains(line, "y") {
		return true
	}
	return false
}

func getCount(label string) int {
	var line string
	fmt.Printf("%s", label)
	fmt.Scanln(&line)
	count, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		panic(err)
	}
	return count
}

var RULES = [...]string{
	"HAS_COLOUR(white)",
	"HAS_COLOUR(black)",
	"HAS_COLOUR(yellow)",
	"HAS_COLOUR(orange)",
	"HAS_COLOUR(red)",
	"HAS_COLOUR(blue)",
	"HAS_TRIANGLE_LEFT()",
	"HAS_COLOUR(green)",
	"HAS_STRIPES(horizontal)",
	"HAS_STRIPES(vertical)",
	"HAS_STRIPES(both-h-and-v)",
	"HAS_STRIPES(diagonal)",
	"HAS_STRIPES_PROPERTY(all_different_colours)",
	"HAS_STRIPES_PROPERTY(alternating_two_colours)",
	"HAS_STRIPES_PROPERTY(uneven) [including borders]",
	"HAS_STRIPES_COUNT(horizontal, 3+)",
	"HAS_STRIPES_COUNT(horizontal, 4+)",
	"HAS_STRIPES_COUNT(vertical, 3+)",
	"HAS_CROSS(any)",
	"HAS_CROSS(diagonal)",
	"HAS_CROSS(greek)",
	"HAS_CROSS(nordic)",
	"HAS_CROSS_PROPERTY(any, border)",
	"HAS_STRIPES_COUNT(horizontal, 6+)",
	"HAS_MIDDLE_CIRCLE()",
	"HAS_SYMMETRY_AXES(any)",
	"HAS_SYMMETRY_AXES(horizontal)",
	"HAS_SYMMETRY_AXES(vertical)",
	"HAS_SYMMETRY_AXES(both)",
	"HAS_FEATURE(star)",
	"HAS_FEATURE_COLOUR(star, white)",
	"HAS_FEATURE_COLOUR(star, yellow)",
	"HAS_FEATURE_COUNT(star, 2+)",
	"HAS_FEATURE_COUNT(star, 3+)",
	"HAS_FEATURE_COUNT(star, 4+)",
	"HAS_FEATURE_COUNT(star, 5+)",
	"HAS_FEATURE_COUNT(star, 6+)",
	"HAS_FEATURE_COUNT(star, 10+)",
	"HAS_FEATURE(sun)",
	"HAS_FEATURE_PROPERTY(sun, half-disk)",
	"HAS_FEATURE(crescent-moon)",
	"HAS_FEATURE_PROPERTY(crescent-moon, pointing-sideways)",
	"HAS_FEATURES(crescent-moon && star)",
	"HAS_FEATURE(writing)",
	"HAS_FEATURE_PROPERTY(writing, arabic)",
	"HAS_FEATURE_PROPERTY(writing, non-its-name-latin)",
	"HAS_FEATURE_PROPERTY(writing, not-its-name-eng)",
	"HAS_FEATURE_PROPERTY(writing, not-its-name-spa-por)",
	"HAS_FEATURE(plant)",
	"HAS_FEATURE_PROPERTY(plant, tree)",
	"HAS_FEATURE(weapon)",
	"HAS_FEATURE_PROPERTY(spear || lance)",
	"HAS_FEATURE(tool)",
	"HAS_FEATURE(beast)",
	"HAS_FEATURE_PROPERTY(beast, mythical)",
	"HAS_FEATURE_PROPERTY(beast, bird)",
	"HAS_FEATURE(uk-flag)",
	"HAS_FEATURE(southern-cross)",
	"HAS_FEATURE(naval-motif)",
	"HAS_FEATURE(coat-of-arms)",
	"HAS_FEATURE(crown)",
	"HAS_FEATURE(human-figures)",
	"HAS_FEATURE(building)",
	"HAS_FEATURE(waves)"}
