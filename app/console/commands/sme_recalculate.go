package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// SmeRecalculate is the artisan command for recalculating SME classifications and formalisation scores
type SmeRecalculate struct {
}

// Signature The name and signature of the console command.
func (receiver *SmeRecalculate) Signature() string {
	return "sme:recalculate"
}

// Description The console command description.
func (receiver *SmeRecalculate) Description() string {
	return "Recalculate formalisation scores and classifications for all SMEs"
}

// Extend The console command extend.
func (receiver *SmeRecalculate) Extend() command.Extend {
	return command.Extend{
		Category: "sme",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "classification-only",
				Aliases: []string{"c"},
				Usage:   "Only recalculate classifications (skip formalisation scores)",
			},
			&command.BoolFlag{
				Name:    "score-only",
				Aliases: []string{"s"},
				Usage:   "Only recalculate formalisation scores (skip classifications)",
			},
			&command.IntFlag{
				Name:    "sme-id",
				Aliases: []string{"i"},
				Usage:   "Recalculate for a specific SME ID only",
			},
			&command.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Usage:   "Show what would be recalculated without making changes",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *SmeRecalculate) Handle(ctx console.Context) error {
	classificationOnly := ctx.OptionBool("classification-only")
	scoreOnly := ctx.OptionBool("score-only")
	smeID := ctx.OptionInt("sme-id")
	dryRun := ctx.OptionBool("dry-run")

	if classificationOnly && scoreOnly {
		ctx.Error("Cannot use --classification-only and --score-only together")
		return fmt.Errorf("conflicting options")
	}

	if dryRun {
		ctx.Warning("DRY RUN MODE - No changes will be made")
	}

	smeService := services.NewSmeService()

	// Single SME mode
	if smeID > 0 {
		return receiver.recalculateSingle(ctx, smeService, uint(smeID), classificationOnly, scoreOnly, dryRun)
	}

	// All SMEs mode
	return receiver.recalculateAll(ctx, smeService, classificationOnly, scoreOnly, dryRun)
}

func (receiver *SmeRecalculate) recalculateSingle(ctx console.Context, smeService *services.SmeService, smeID uint, classificationOnly, scoreOnly, dryRun bool) error {
	// Verify SME exists
	var sme models.Sme
	err := facades.Orm().Query().
		With("BusinessFormalisation").
		With("BusinessEmployeeSummary").
		With("PrimaryBusinessOwner").
		With("AdditionalBusinessMembers").
		Where("id = ?", smeID).
		First(&sme)

	if err != nil || sme.ID == 0 {
		ctx.Error(fmt.Sprintf("SME with ID %d not found", smeID))
		return fmt.Errorf("SME not found")
	}

	ctx.Info(fmt.Sprintf("Processing SME: %s (ID: %d)", sme.Name, sme.ID))
	ctx.Info(fmt.Sprintf("  Current classification: %s", sme.Classification))

	if sme.BusinessFormalisation != nil {
		ctx.Info(fmt.Sprintf("  Current formalisation score: %d", sme.BusinessFormalisation.FormalisationScore))
	}

	if dryRun {
		ctx.Info("  [DRY RUN] Would recalculate scores and classification")
		return nil
	}

	var newScore int
	var newClassification string

	// Recalculate formalisation score
	if !classificationOnly {
		newScore, err = smeService.CalculateFormalisationScore(smeID)
		if err != nil {
			ctx.Error(fmt.Sprintf("  Failed to calculate formalisation score: %v", err))
		} else {
			ctx.Success(fmt.Sprintf("  New formalisation score: %d", newScore))
		}
	}

	// Recalculate classification
	if !scoreOnly {
		newClassification, err = smeService.CalculateClassification(smeID)
		if err != nil {
			ctx.Error(fmt.Sprintf("  Failed to calculate classification: %v", err))
		} else {
			ctx.Success(fmt.Sprintf("  New classification: %s", newClassification))
		}
	}

	return nil
}

func (receiver *SmeRecalculate) recalculateAll(ctx console.Context, smeService *services.SmeService, classificationOnly, scoreOnly, dryRun bool) error {
	// Get all SME IDs
	var smes []models.Sme
	err := facades.Orm().Query().
		Model(&models.Sme{}).
		Select("id", "name", "classification").
		Find(&smes)

	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to fetch SMEs: %v", err))
		return err
	}

	totalCount := len(smes)
	ctx.Info(fmt.Sprintf("Found %d SMEs to process", totalCount))

	if dryRun {
		ctx.Info("[DRY RUN] Would recalculate for all SMEs")
		return nil
	}

	var scoreSuccessCount, scoreErrorCount int
	var classificationSuccessCount, classificationErrorCount int
	var classificationChanges []string

	for i, sme := range smes {
		// Progress indicator every 100 SMEs
		if (i+1)%100 == 0 || i == 0 {
			ctx.Info(fmt.Sprintf("Processing %d/%d...", i+1, totalCount))
		}

		oldClassification := sme.Classification

		// Recalculate formalisation score
		if !classificationOnly {
			_, err := smeService.CalculateFormalisationScore(sme.ID)
			if err != nil {
				scoreErrorCount++
			} else {
				scoreSuccessCount++
			}
		}

		// Recalculate classification
		if !scoreOnly {
			newClassification, err := smeService.CalculateClassification(sme.ID)
			if err != nil {
				classificationErrorCount++
			} else {
				classificationSuccessCount++
				if oldClassification != newClassification {
					classificationChanges = append(classificationChanges,
						fmt.Sprintf("  %s (ID: %d): %s -> %s", sme.Name, sme.ID, oldClassification, newClassification))
				}
			}
		}
	}

	// Print summary
	ctx.NewLine()
	ctx.Info("========== RECALCULATION SUMMARY ==========")
	ctx.Info(fmt.Sprintf("Total SMEs: %d", totalCount))

	if !classificationOnly {
		ctx.Success(fmt.Sprintf("Formalisation scores updated: %d", scoreSuccessCount))
		if scoreErrorCount > 0 {
			ctx.Error(fmt.Sprintf("Formalisation score errors: %d", scoreErrorCount))
		}
	}

	if !scoreOnly {
		ctx.Success(fmt.Sprintf("Classifications updated: %d", classificationSuccessCount))
		if classificationErrorCount > 0 {
			ctx.Error(fmt.Sprintf("Classification errors: %d", classificationErrorCount))
		}

		if len(classificationChanges) > 0 {
			ctx.NewLine()
			ctx.Info(fmt.Sprintf("Classification changes (%d):", len(classificationChanges)))
			for _, change := range classificationChanges {
				ctx.Info(change)
			}
		} else {
			ctx.Info("No classification changes detected")
		}
	}

	ctx.Info("============================================")

	return nil
}
