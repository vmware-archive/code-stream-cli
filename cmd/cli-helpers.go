/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// PrettyPrint prints interfaces
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// PrintTable prints an array of objects with table headers
func PrintTable(objects []interface{}, headers []string) {
	// Print result table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, object := range objects {
		t := reflect.TypeOf(object)
		fmt.Println(t)
		// var o t
		// mapstructure.Decode(object, &t)
		// var values []string
		// for _, header := range headers {
		// 	append(values, object.(*t).header)
		// }
		// 	table.Append([]string{c.ID, c.Name, c.Project})
	}
	table.Render()
}

func getYamlFilePaths(importPath string) []string {
	var yamlFiles []string
	// Read importPath
	stat, err := os.Stat(importPath)
	if err == nil && stat.IsDir() {
		// log.Debugln("importPath is a directory")
		files, err := ioutil.ReadDir(importPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if strings.Contains(f.Name(), ".yaml") || strings.Contains(f.Name(), ".yml") {
				yamlFiles = append(yamlFiles, filepath.Join(importPath, f.Name()))
			}
		}
	} else {
		// log.Debugln("importPath is a file")
		yamlFiles = append(yamlFiles, importPath)
	}
	return yamlFiles
}

func removeDuplicateStrings(elements []string) []string {
	encountered := map[string]bool{}
	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}
	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
