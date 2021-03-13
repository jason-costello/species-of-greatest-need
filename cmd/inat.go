// Package cmd /*
package cmd

import (
	"fmt"
	"sogn/pkg/inaturalist"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var client *inaturalist.INaturalist

// inatCmd represents the inat command
var inatCmd = &cobra.Command{
	Use:   "inat",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("inat called")
	},
}

var taxaCmd = &cobra.Command{
	Use: "taxa",
	Short: "Return details for a taxa id",
	Long: "Return details for a taxa id",
	Run: func(cmd *cobra.Command, args []string){

		var ids []int
		for _, j := range args{
				i, _ := strconv.Atoi(j)
				ids = append(ids, i)
				}

		d, err := client.GetTaxonDetails(ids...)
		if err != nil{
			fmt.Println("error requesting details for taxa id.", err.Error())
		}

		for _, r := range d.Results{

			fmt.Printf("%s: %s\tCommon Name: %s\tObservations: %d\tTaxaID: %d \n", strings.Title(r.Rank), r.Name, r.PreferredCommonName,r.ObservationsCount, r.ID)




		}

	},
}

var observationCmd = &cobra.Command{
	Use: "observation",
	Short: "Return group of observations",
	Long: "Return group of observations",
	Run: func(cmd *cobra.Command,  args []string){

		op := inaturalist.ObservationParameters{}


		for _, s := range args{

			ss := strings.Split(s, "=")
			if len(ss) > 1 {
				op[ss[0]] = ss[1]
			}

		}

		d, err := client.Observations(op)
		if err != nil{
			fmt.Println("error requesting details for taxa id.", err.Error())
		}

		for _, r := range d.Results{

			fmt.Printf("obs id: %d   location: %s  quality: %s\n", r.ID, r.Location,  r.QualityGrade)


		}

	},
}

func init() {
	client = inaturalist.NewClient()

	rootCmd.AddCommand(inatCmd)
	inatCmd.AddCommand(taxaCmd)
	inatCmd.AddCommand(observationCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inatCmd.PersistentFlags().String("taxa", "", "return taxa details for id")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
