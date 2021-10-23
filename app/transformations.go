package app

import "meanos.io/console/review"

func (a *Application) MakeExported() ExportedApplication {
	var ea ExportedApplication
	ea.Name = a.Name
	ea.Publisher = a.Publisher.DisplayName
	ea.PublisherVerified = a.Publisher.Verified
	ea.Downloads = a.Downloads
	ea.Version = a.Version.CurrentVersion.Version
	ea.ReleaseNotes = a.Version.CurrentVersion.ReleaseNotes
	ea.Category = a.Category
	ea.Reviews = review.ExportSlice(a.Reviews)
	ea.Rating = a.Rating

	ea.IconURL = a.IconURL
	ea.CoverURL = a.CoverURL
	ea.PrettyPreviewURL = a.PrettyPreviewURL
	ea.ScreenshotsURLs = a.ScreenshotsURLs

	ea.Monthly, ea.Yearly, ea.Once, ea.Free, ea.Price = a.PaymentType.ExportFields()

	if ea.Free {
		ea.PackageURL = a.Version.CurrentVersion.PackageURL
	} else {
		ea.PackageURL = ""
	}

	return ea
}

func ExportSlice(apps []Application) []ExportedApplication {
	var result []ExportedApplication
	for _, a := range apps {
		result = append(result, a.MakeExported())
	}

	return result
}
