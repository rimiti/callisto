package models

import (
	"github.com/graphql-go/graphql"
)

type Release struct {
	Id        int    `json:"id",db:"id"`
	ProjectId int    `json:"projectId",db:"projectId"`
	Version   string `json:"version",db:"version"`
	Deleted   int    `json:"deleted",db:"deleted"`
}

// define custom GraphQL ObjectType `ReleaseType` for our Golang struct `Release`
// Note that
// - the fields in our ReleaseType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var ReleaseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Release",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"projectId": &graphql.Field{
			Type: graphql.Int,
		},
		"version": &graphql.Field{
			Type: graphql.String,
		},
		"deleted": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

func AllReleases() ([]*Release, error) {
	rows, err := db.Query(`SELECT * FROM releases`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rels := make([]*Release, 0)

	for rows.Next() {
		rel := new(Release)
		err := rows.Scan(&rel.Id, &rel.ProjectId, &rel.Version, &rel.Deleted)
		if err != nil {
			return nil, err
		}
		rels = append(rels, rel)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rels, nil
}

func ReleasesForProject(id int) ([]*Release, error) {
	rows, err := db.Query(`SELECT * FROM releases WHERE projectId=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	releases := make([]*Release, 0)

	for rows.Next() {
		release := new(Release)
		err := rows.Scan(&release.Id, &release.ProjectId, &release.Version, &release.Deleted)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return releases, nil
}

func NewRelease(rel Release) (error) {
	stmt, err := db.Prepare("INSERT INTO releases(version, projectId, deleted) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rel.Version, rel.ProjectId, rel.Deleted)
	return err
}
