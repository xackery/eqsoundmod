package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"bytes"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/wav"
	"github.com/xackery/quail/pfs"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed: ", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 4 {
		return fmt.Errorf("usage: eqsoundmod <eqpath> <wav> <volumeadj>")
	}
	sndFiles := []string{
		"snd1.pfs",
		"snd2.pfs",
		"snd3.pfs",
		"snd4.pfs",
		"snd5.pfs",
		"snd6.pfs",
		"snd7.pfs",
		"snd8.pfs",
		"snd9.pfs",
		"snd10.pfs",
		"snd11.pfs",
		"snd12.pfs",
		"snd13.pfs",
		"snd14.pfs",
		"snd15.pfs",
		"snd16.pfs",
		"snd17.pfs",
	}
	path := os.Args[1]
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	targetFile := os.Args[2]

	volumeStr := os.Args[3]
	volumeInt, err := strconv.Atoi(volumeStr)
	if err != nil {
		return fmt.Errorf("invalid volume %s, must be range of 0 to 100", volumeStr)
	}
	if volumeInt < 0 || volumeInt > 100 {
		return fmt.Errorf("invalid volume %s, must be range of 0 to 100", volumeStr)
	}
	// convert int 0 to 100 to float64 0.0 to 1.0
	volumeFloat := float64(volumeInt) / 100.0

	fmt.Printf("Volume: %f\n", volumeFloat)

	volumeFloat = 0.7 + (volumeFloat * 0.3)

	outPath := filepath.Join(path, "sounds")
	_, err = os.Stat(outPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("output path %s does not exist", outPath)
		}
		return fmt.Errorf("stat %s: %w", outPath, err)
	}

	for _, sndFile := range sndFiles {
		sndPath := filepath.Join(path, sndFile)
		archive, err := pfs.NewFile(sndPath)
		if err != nil {
			return fmt.Errorf("open %s: %w", sndPath, err)
		}
		defer archive.Close()

		for _, file := range archive.Files() {
			if file.Name() != targetFile {
				continue
			}

			err = lowerVolume(bytes.NewReader(file.Data()), filepath.Join(outPath, targetFile), volumeFloat)
			if err != nil {
				return fmt.Errorf("lower volume: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("file %s not found in any snd files", targetFile)
}

func lowerVolume(r *bytes.Reader, outputFile string, volumeFactor float64) error {

	// Decode WAV file into streamer and format
	streamer, format, err := wav.Decode(r)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	defer streamer.Close()
	volumeDB := 20 * math.Log10(volumeFactor)
	// Apply volume reduction
	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   volumeDB,
	}

	// Open output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer outFile.Close()

	// Encode the volume-adjusted stream to WAV
	err = wav.Encode(outFile, volume, format)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}
