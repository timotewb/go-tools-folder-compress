package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ncruces/zenity"
)

func main(){
	inDir, err := zenity.SelectFile(
		zenity.Filename(""),
		zenity.Directory(),
		zenity.DisallowEmpty(),
		zenity.Title("Select input directory."),
	)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// find folders
	files, err := os.ReadDir(inDir)
	if err != nil{
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println("")
		fmt.Println("----------------------------------------------------------------------------------------")
		fmt.Println(file.Name())
		fmt.Println("----------------------------------------------------------------------------------------")

		if file.IsDir(){
			dir := filepath.Join(inDir,file.Name())
			zipFile, err := os.Create(dir + ".zip")
			if err != nil {
				log.Fatal("os.Create()"+err.Error())
			}
			defer zipFile.Close()

			zipWriter := zip.NewWriter(zipFile)
			defer zipWriter.Close()

			fmt.Println("dir:",dir)
			fmt.Println("zipFile:",dir + ".zip")

			err = filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatal("filepath.Walk()"+err.Error())
					return err
				}

				// Create a local file header
				header, err := zip.FileInfoHeader(info)
				if err != nil {
					log.Fatal("zip.FileInfoHeader()"+err.Error())
					return err
				}

				// Set compression
				header.Method = zip.Deflate

				// Set relative path of a file as the header name
				header.Name, err = filepath.Rel(filepath.Dir(dir), filePath)
				if err != nil {
					log.Fatal("filepath.Rel()"+err.Error())
					return err
				}
				if info.IsDir() {
					header.Name += "/"
				}

				// Create writer for the file header and save content of the file
				headerWriter, err := zipWriter.CreateHeader(header)
				if err != nil {
					log.Fatal("zipWriter.CreateHeader()"+err.Error())
					return err
				}
				if info.IsDir() {
					return nil
				}
				if len(filepath.Base(filePath)) > 0 && filepath.Base(filePath)[0] != '.' {
					f, err := os.Open(filePath)
					if err != nil {
						log.Fatal("os.Open()"+err.Error())
						return err
					}
					defer f.Close()

					_, err = io.Copy(headerWriter, f)
					if err != nil {
						log.Fatal("io.Copy()"+err.Error())
						return err
					}
					return nil
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	zenity.Info("Folders compressed!",
		zenity.Title("Complete"),
		zenity.InfoIcon,
	)
}