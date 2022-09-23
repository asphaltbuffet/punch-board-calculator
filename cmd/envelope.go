package cmd

import (
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/punch-board-calculator/pkg/calculate"
)

const envelopeCommandLongDesc = "LONG DESCRIPTION GOES HERE."

// NewEnvelopeCommand returns a new envelope command.
func NewEnvelopeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "envelope",
		Short: "calculate punch positions for an envelope",
		Long:  envelopeCommandLongDesc,
		Run:   RunEnvelopeCmd,
	}

	cmd.Flags().Float64P("length", "l", 0, "length of envelope")
	cmd.Flags().Float64P("width", "w", 0, "width of envelope")
	cmd.Flags().Bool("loose", false, "loose envelope")
	cmd.Flags().Bool("mini", false, "mini punch board")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewEnvelopeCommand())
}

// RunEnvelopeCmd is the entrypoint for the envelope command.
func RunEnvelopeCmd(cmd *cobra.Command, args []string) {
	length, _ := cmd.Flags().GetFloat64("length")
	width, _ := cmd.Flags().GetFloat64("width")

	isLoose, _ := cmd.Flags().GetBool("loose")
	isMini, _ := cmd.Flags().GetBool("mini")

	cmd.Printf("Content (length x width): %0.2f x %0.2f\n", length, width)

	paperSize, punchLocation, err := calculate.CalculateEnvelope(length, width, isLoose, isMini)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Printf("Paper size: %0.1f\n", paperSize)
	cmd.Printf("Punch location: %0.1f\n", punchLocation)
}
