package readmodel

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestNamePermissionsAdministratorShortCircuits(t *testing.T) {
	got := namePermissions(discordgo.PermissionAdministrator | discordgo.PermissionKickMembers)
	if !reflect.DeepEqual(got, []string{"Administrator"}) {
		t.Fatalf("expected only Administrator, got %v", got)
	}
}

func TestNamePermissionsListsGrantedOnly(t *testing.T) {
	got := namePermissions(discordgo.PermissionKickMembers | discordgo.PermissionBanMembers)
	want := map[string]bool{"Kick Members": true, "Ban Members": true}
	if len(got) != 2 {
		t.Fatalf("expected 2 permissions, got %v", got)
	}
	for _, name := range got {
		if !want[name] {
			t.Fatalf("unexpected permission %q in %v", name, got)
		}
	}
}

func TestNamePermissionsEmpty(t *testing.T) {
	got := namePermissions(0)
	if len(got) != 0 {
		t.Fatalf("expected no permissions, got %v", got)
	}
}

func TestRolePermissionsCombinesEveryoneAndMemberRoles(t *testing.T) {
	everyone := &discordgo.Role{ID: "guild1", Permissions: discordgo.PermissionViewChannel}
	roles := []*discordgo.Role{
		everyone,
		{ID: "role1", Permissions: discordgo.PermissionSendMessages},
		{ID: "role2", Permissions: discordgo.PermissionKickMembers},
	}

	got := rolePermissions(everyone, roles, []string{"role1"})
	want := int64(discordgo.PermissionViewChannel | discordgo.PermissionSendMessages)
	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

func TestRolePermissionsAdministratorImpliesAll(t *testing.T) {
	roles := []*discordgo.Role{
		{ID: "admin-role", Permissions: discordgo.PermissionAdministrator},
	}
	got := rolePermissions(nil, roles, []string{"admin-role"})
	if got&discordgo.PermissionKickMembers == 0 {
		t.Fatal("expected administrator to imply all other permissions")
	}
}

func TestTopRolePositionPicksHighest(t *testing.T) {
	roles := []*discordgo.Role{
		{ID: "a", Position: 1},
		{ID: "b", Position: 5},
		{ID: "c", Position: 3},
	}
	got := topRolePosition(roles, []string{"a", "b", "c"})
	if got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestTopRolePositionNoMatchingRoles(t *testing.T) {
	roles := []*discordgo.Role{{ID: "a", Position: 5}}
	got := topRolePosition(roles, []string{"z"})
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestFindEveryoneRole(t *testing.T) {
	roles := []*discordgo.Role{
		{ID: "guild1"},
		{ID: "role1"},
	}
	got := findEveryoneRole(roles, "guild1")
	if got == nil || got.ID != "guild1" {
		t.Fatalf("expected to find the @everyone role, got %v", got)
	}

	if findEveryoneRole(roles, "missing") != nil {
		t.Fatal("expected nil when no role matches the guild ID")
	}
}
