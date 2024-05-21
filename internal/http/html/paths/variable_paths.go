// Code generated by "go generate"; DO NOT EDIT.

package paths

import "fmt"

func Variables(workspace string) string {
    return fmt.Sprintf("/app/workspaces/%s/variables", workspace)
}

func CreateVariable(workspace string) string {
    return fmt.Sprintf("/app/workspaces/%s/variables/create", workspace)
}

func NewVariable(workspace string) string {
    return fmt.Sprintf("/app/workspaces/%s/variables/new", workspace)
}

func Variable(variable string) string {
    return fmt.Sprintf("/app/variables/%s", variable)
}

func EditVariable(variable string) string {
    return fmt.Sprintf("/app/variables/%s/edit", variable)
}

func UpdateVariable(variable string) string {
    return fmt.Sprintf("/app/variables/%s/update", variable)
}

func DeleteVariable(variable string) string {
    return fmt.Sprintf("/app/variables/%s/delete", variable)
}

