package domain_models

type SavePicParams struct {
	File	string
}

type ConfirmSaveParams struct {
	Name		string
	UserID		uint
}

type File	struct {
	Name	string
	Link	string
}