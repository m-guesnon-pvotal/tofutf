// Code generated by "go generate"; DO NOT EDIT.

package paths

import "fmt"

func Teams(organization string) string {
    return fmt.Sprintf("/app/organizations/%s/teams", organization)
}

func CreateTeam(organization string) string {
    return fmt.Sprintf("/app/organizations/%s/teams/create", organization)
}

func NewTeam(organization string) string {
    return fmt.Sprintf("/app/organizations/%s/teams/new", organization)
}

func Team(team string) string {
    return fmt.Sprintf("/app/teams/%s", team)
}

func EditTeam(team string) string {
    return fmt.Sprintf("/app/teams/%s/edit", team)
}

func UpdateTeam(team string) string {
    return fmt.Sprintf("/app/teams/%s/update", team)
}

func DeleteTeam(team string) string {
    return fmt.Sprintf("/app/teams/%s/delete", team)
}

func AddMemberTeam(team string) string {
    return fmt.Sprintf("/app/teams/%s/add-member", team)
}

func RemoveMemberTeam(team string) string {
    return fmt.Sprintf("/app/teams/%s/remove-member", team)
}

