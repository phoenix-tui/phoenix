// Package main demonstrates how to create dialog boxes with Phoenix layout.
package main

import (
	"fmt"

	layout "github.com/phoenix-tui/phoenix/layout/api"
	"github.com/phoenix-tui/phoenix/layout/domain/value"
)

func main() {
	fmt.Println("=== Dialog Box Examples ===")

	// Example 1: Simple confirmation dialog
	fmt.Println("1. Simple Confirmation Dialog:")
	confirm := layout.NewBox("Are you sure?").
		PaddingAll(2).
		Border().
		AlignCenter().
		Render()
	fmt.Println(confirm)
	fmt.Println()

	// Example 2: Multi-line dialog with buttons
	fmt.Println("2. Dialog with Buttons:")
	dialog := layout.NewBox("Delete this file?\n\n[Yes] [No]").
		PaddingVH(2, 3).
		Border().
		Render()
	fmt.Println(dialog)
	fmt.Println()

	// Example 3: Warning dialog
	fmt.Println("3. Warning Dialog:")
	warning := layout.NewBox("⚠️  Warning!\n\nThis action cannot be undone.\n\n[Continue] [Cancel]").
		PaddingAll(2).
		Border().
		MarginAll(1).
		Render()
	fmt.Println(warning)
	fmt.Println()

	// Example 4: Info dialog with title
	fmt.Println("4. Info Dialog:")
	info := layout.NewBox("ℹ️  Information\n\nThe process completed successfully.\nYou can now close this dialog.\n\n[OK]").
		PaddingVH(1, 3).
		Border().
		Render()
	fmt.Println(info)
	fmt.Println()

	// Example 5: Error dialog
	fmt.Println("5. Error Dialog:")
	errorBox := layout.NewBox("❌ Error\n\nFailed to save file:\nPermission denied\n\n[Retry] [Cancel]").
		PaddingAll(2).
		Border().
		Render()
	fmt.Println(errorBox)
	fmt.Println()

	// Example 6: Dialog with layout positioning
	fmt.Println("6. Dialog Positioning:")
	positionedDialog := layout.NewBox("Centered Dialog\n\nThis dialog is positioned\nin the center of an 80x24 terminal.").
		PaddingAll(2).
		Border()

	pos := positionedDialog.Layout(80, 24)
	fmt.Printf("Dialog would be rendered at position (%d, %d)\n", pos.X(), pos.Y())
	fmt.Println(positionedDialog.Render())
	fmt.Println()

	// Example 7: Form dialog
	fmt.Println("7. Form Dialog:")
	form := layout.NewBox("Login\n\nUsername: ___________\nPassword: ___________\n\n[Login] [Cancel]").
		PaddingVH(2, 4).
		Border().
		Render()
	fmt.Println(form)
	fmt.Println()

	// Example 8: Progress dialog
	fmt.Println("8. Progress Dialog:")
	progress := layout.NewBox("Processing...\n\n[████████░░] 80%\n\nPlease wait...").
		PaddingAll(2).
		Border().
		Render()
	fmt.Println(progress)
	fmt.Println()

	// Example 9: Dialog with custom alignment
	fmt.Println("9. Right-aligned Dialog:")
	rightDialog := layout.NewBox("Notification\n\nYou have 3 new messages.").
		PaddingAll(1).
		Border().
		Align(value.AlignRight, value.AlignTop)

	rightPos := rightDialog.Layout(80, 24)
	fmt.Printf("Would be at position (%d, %d)\n", rightPos.X(), rightPos.Y())
	fmt.Println(rightDialog.Render())
}
