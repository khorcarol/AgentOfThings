//go:generate fyne bundle -o bundled.go --package frontend Inter_24pt-Bold.ttf

package frontend

import (
	// "time"
	"image"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/api/interests"
	"github.com/khorcarol/AgentOfThings/internal/middle"
)

// Custom colors
var (
	primaryColor    = color.NRGBA{R: 0, G: 120, B: 215, A: 255}   // Blue
	backgroundColor = color.NRGBA{R: 245, G: 245, B: 245, A: 255} // Off-white
	listItemColor   = color.White
)

// Custom theme implementation
type customTheme struct {
	fyne.Theme
	regularFont     fyne.Resource
	primaryColor    color.Color
	backgroundColor color.Color
}

func (c *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	if !style.Bold && !style.Italic {
		return c.regularFont
	}
	return c.Theme.Font(style)
}

func (c *customTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return c.primaryColor
	case theme.ColorNameBackground:
		return c.backgroundColor
	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

var (
	friendsList    *widget.List
	usersList      *widget.List
	currentFriends []api.Friend
	currentUsers   []api.User
)

func createFriendsUI() fyne.CanvasObject {
	friendsList = widget.NewList(
		func() int { return len(currentFriends) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				canvas.NewImageFromResource(nil),
				container.NewVBox(
					widget.NewLabel("Name"),
					widget.NewLabel("Common Interests"),
				),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			friend := currentFriends[i]
			container := o.(*fyne.Container)
			img := container.Objects[0].(*canvas.Image)
			details := container.Objects[1].(*fyne.Container)

			img.Resource = nil
			img.SetMinSize(fyne.NewSize(48, 48))

			nameLabel := details.Objects[0].(*widget.Label)
			nameLabel.SetText(friend.Name)

			interestsLabel := details.Objects[1].(*widget.Label)
			interestsLabel.SetText(formatInterests(middle.CommonInterests(friend.User)))
		},
	)

	return container.NewBorder(nil, nil, nil, nil, friendsList)
}

func onRefreshFriends(friends []api.Friend) {
	currentFriends = friends
	if friendsList != nil {
		friendsList.Refresh()
	}
}

func onRefreshUsers(users []api.User) {
	currentUsers = users
	if usersList != nil {
		usersList.Refresh()
	}
}

func formatInterests(interests []api.Interest) string {
	result := ""
	for _, interest := range interests {
		result += interest.Description + "\n"
	}
	return result
}

func createUsersUI(myWindow fyne.Window) fyne.CanvasObject {

	usersList = widget.NewList(
		func() int { return len(currentUsers) },
		func() fyne.CanvasObject {
			image := &canvas.Image{}
			image.SetMinSize(fyne.Size{Width: 200, Height: 200})
			return container.NewVBox(container.NewHBox(
				widget.NewLabel("User ID"),
				layout.NewSpacer(),
				widget.NewButton("Learn More", nil),
			),
				image,
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := currentUsers[i]
			vertContainer := o.(*fyne.Container)
			container := vertContainer.Objects[0].(*fyne.Container)

			button := container.Objects[2].(*widget.Button)
			button.OnTapped = func() {
				middle.Seen(user.UserID)
				showUserDetailsDialog(user, myWindow)
			}

			image := vertContainer.Objects[1].(*canvas.Image)

			image.Resource, _ = fyne.LoadResourceFromURLString("https://go.dev/blog/go-brand/Go-Logo/JPG/Go-Logo_Blue.jpg")
			image.FillMode = canvas.ImageFillContain
		},
	)
	return container.NewBorder(nil, nil, nil, nil, usersList)
}

func showUserDetailsDialog(user api.User, parent fyne.Window) {
	content := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("Common Interests:"),
	)

	for _, interest := range middle.CommonInterests(user) {
		content.Add(widget.NewLabel("- " + interests.String(interest.Category) + ": " + interest.Description))
	}

	var userDetailsDialog dialog.Dialog

	sendBtn := widget.NewButton("Send Friend Request", func() {
		middle.SendFriendRequest(user.UserID)
		dialog.ShowInformation("Request Sent", "Friend request sent!", parent)
	})
	closeBtn := widget.NewButton("Close", func() {
		userDetailsDialog.Hide()
	})

	buttons := container.NewHBox(
		sendBtn,
		closeBtn,
	)
	content.Add(buttons)
	userDetailsDialog = dialog.NewCustomWithoutButtons("User Details",
		content,
		parent,
	)

	userDetailsDialog.Show()
}

func Main() {
	regularFont := resourceInter24ptBoldTtf

	myApp := app.New()
	myWindow := myApp.NewWindow("Agent of Friends")
	myWindow.Resize(fyne.NewSize(800, 600))
	myApp.Settings().SetTheme(&customTheme{
		Theme:           theme.DefaultTheme(),
		regularFont:     regularFont,
		primaryColor:    primaryColor,
		backgroundColor: backgroundColor,
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Peers", createUsersUI(myWindow)),
		container.NewTabItem("Friends", createFriendsUI()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)

	// temp, until real friends used
	reader, err := os.Open("assets/blank-profile.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Print(err)
	}

	onRefreshUsers([]api.User{
		{
			UserID: api.ID{Address: uuid.Nil},
			Interests: []api.Interest{
				{Category: 1, Description: "Basketball"},
				{Category: 2, Description: "Jazz"},
			},
			Seen: false,
		},
		{
			UserID: api.ID{Address: uuid.Nil},
			Interests: []api.Interest{
				{Category: 1, Description: "Basketball"},
				{Category: 2, Description: "Jazz"},
			},
			Seen: false,
		},
	})

	onRefreshFriends([]api.Friend{
		{
			User: api.User{
				UserID:    api.ID{Address: uuid.Nil},
				Interests: []api.Interest{{Category: 1, Description: "description"}},
			},
			Name:  "Friend",
			Photo: img,
		},
	})
	myWindow.ShowAndRun()
}
