package review

type Review struct {
	Name   string  `bson:"name"    json:"name"`
	AppId  string  `bson:"app_id"  json:"app_id"`
	Text   string  `bson:"text"    json:"text"`
	Rating float64 `bson:"rating"  json:"rating"`
	UID    float64 `bson:"uid"     json:"uid"`
}

type ExportedReview struct {
	Name   string  `json:"name"`
	Text   string  `json:"text"`
	Rating float64 `json:"rating"`
}

func (r *Review) MakeExported() ExportedReview {
	var er ExportedReview
	er.Name = r.Name
	er.Text = r.Text
	er.Rating = r.Rating
	return er
}

func ExportSlice(rs []Review) []ExportedReview {
	var res []ExportedReview
	for _, r := range rs {
		res = append(res, r.MakeExported())
	}
	return res
}
