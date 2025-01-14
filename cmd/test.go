package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bata94/reqlab/pkgs/apiview/openapi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test cmd",
	Long:  "Text cmd",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Test cmd")

		jsonFile, err := os.Open("tmp/swagger.json")

		if err != nil {
			log.Error(err)
		}

		defer jsonFile.Close()
		log.Info("Successfully Opened JsonFile")

		openapiCollection := openapi.OpenAPI{}
		jsonDecoder := json.NewDecoder(jsonFile)

		err = jsonDecoder.Decode(&openapiCollection)
		if err != nil {
			log.Fatal(err)
		}

		// log.Info(json.Marshal(openapiCollection))

		// Open a file for writing (create it if it doesn't exist)
		file, err := os.Create(fmt.Sprintf("tmp/test%s.json", time.Now().Format("2006-01-02_15:04:05")))
		if err != nil {
			log.Error("Error creating file:", err)
			return
		}
		defer file.Close() // Ensure the file is closed when done

		// Marshal the struct to JSON
		data, err := json.MarshalIndent(openapiCollection, "", "  ") // json.MarshalIndent adds indentation for readability
		if err != nil {
			log.Error("Error marshaling struct to JSON:", err)
			return
		}

		// Write the JSON data to the file
		_, err = file.Write(data)
		if err != nil {
			log.Error("Error writing to file:", err)
			return
		}

		fmt.Println("Struct successfully written")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
