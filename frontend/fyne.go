//go:generate fyne bundle -o bundled.go --package frontend ../assets/Inter_24pt-Bold.ttf
//go:generate fyne bundle -o bundled.go --package frontend --append ../assets/blank-profile.png

package frontend

import (
	// "go/format"

	"image/color"
	"io"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/middle"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

// Custom colors
var (
	primaryColor    = color.NRGBA{R: 0, G: 120, B: 215, A: 255}   // Blue
	backgroundColor = color.NRGBA{R: 245, G: 245, B: 245, A: 255} // Off-white
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

func getImage(interests []api.Interest) *string {
	for i := range interests {
		if interests[i].Image != nil {
			return interests[i].Image
		}
	}
	return nil
}

var (
	friendsList    *widget.List
	usersList      *widget.List
	currentFriends []api.Friend
	currentUsers   []api.User
)

func createFriendsUI() fyne.CanvasObject {
	friendsList = widget.NewList(
		func() int {
			return len(currentFriends)
		},
		func() fyne.CanvasObject {
			// Prepare images.
			profileImage := &canvas.Image{}
			profileImage.SetMinSize(fyne.NewSize(200, 200))
			interestImage := &canvas.Image{}
			interestImage.SetMinSize(fyne.NewSize(200, 200))

			// Build a VBox for text: Name, Contact and Interests.
			leftBox := container.NewVBox(

				canvas.NewText("  Name", color.Black),
				widget.NewLabel("Interests:"),
				widget.NewLabel("Contact:"), // New field for contact info.
			)

			// Build the list item as an HBox that contains:
			// leftBox, a spacer, and the padded images on the right.
			item := container.NewHBox(
				leftBox,
				layout.NewSpacer(),                 // Push images to the right.
				container.NewPadded(profileImage),  // Profile image with padding.
				container.NewPadded(interestImage), // Interest image with padding.
			)
			return item
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			friend := currentFriends[i]
			item := o.(*fyne.Container)

			leftBox := item.Objects[0].(*fyne.Container)
			// leftBox.Objects[0] = Name, [1] = Contact, [2] = Interests.
			nameLabel := leftBox.Objects[0].(*canvas.Text)
			nameLabel.TextSize = 16
			nameLabel.Text = "  " + friend.Name

			interestsLabel := leftBox.Objects[1].(*widget.Label)
			interestsLabel.SetText(formatInterests(friend.User.Interests))

			contactLabel := leftBox.Objects[2].(*widget.Label)
			contactLabel.TextStyle.Italic = true
			contactLabel.SetText(friend.Contact)

			profileImage := item.Objects[3].(*fyne.Container).Objects[0].(*canvas.Image)
			if friend.Photo.Img != nil {
				profileImage.Image = friend.Photo.Img
				profileImage.FillMode = canvas.ImageFillContain
				profileImage.Show()
				profileImage.Refresh()
			} else {
				profileImage.Hide()
			}

			interestsImage := item.Objects[2].(*fyne.Container).Objects[0].(*canvas.Image)
			if imageUrl := getImage(friend.User.Interests); imageUrl != nil {
				// interestImage is the next element, wrapped in a padded container; index [3]
				if resource, err := fyne.LoadResourceFromURLString(*imageUrl); err == nil {
					interestsImage.Resource = resource

					interestsImage.FillMode = canvas.ImageFillContain
					interestsImage.Show()
					interestsImage.Refresh()
				} else {
					interestsImage.Hide()
				}
			} else {
				interestsImage.Hide()
			}
		},
	)

	return container.NewPadded(friendsList)
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

func createUsersUI() fyne.CanvasObject {
	usersList = widget.NewList(
		func() int { return len(currentUsers) },
		func() fyne.CanvasObject {
			image := &canvas.Image{}
			image.SetMinSize(fyne.Size{Width: 200, Height: 200})
			return container.NewVBox(
				canvas.NewText("  Anonymous User", color.Black),
				widget.NewLabel("Interests: "),
				image,
				widget.NewButton("Send Friend Request", nil),
				widget.NewButton("Reject Friend Request", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			user := currentUsers[i]
			container := o.(*fyne.Container)

			nameLabel := container.Objects[0].(*canvas.Text)
			nameLabel.TextSize = 16

			interests_label := container.Objects[1].(*widget.Label)
			interests_label.Text = "Interests: \n" + formatInterests(user.Interests)
			interests_label.Refresh()

			image := container.Objects[2].(*canvas.Image)
			if imageUrl := getImage(user.Interests); imageUrl != nil {
				image.Show()
				image.Resource, _ = fyne.LoadResourceFromURLString(*imageUrl)
				image.FillMode = canvas.ImageFillContain
			} else {
				image.Hide()
			}

			friendButton := container.Objects[3].(*widget.Button)

			friendButton.OnTapped = func() {
				middle.Seen(user.UserID)
				middle.SendFriendRequest(user.UserID, true)
			}

			rejectButton := container.Objects[4].(*widget.Button)

			rejectButton.OnTapped = func() {
				middle.Seen(user.UserID)
				middle.SendFriendRequest(user.UserID, false)
			}
			rejectButton.Hide()
			friendButton.Enable()

			if middle.HasIncomingFriendRequest(user.UserID) {
				friendButton.SetText("Accept Friend Request")
				rejectButton.Show()
			} else if middle.HasOutgoingFriendRequest(user.UserID) {
				friendButton.SetText("Friend Request Sent")
				friendButton.Disable()
			} else {
				friendButton.SetText("Send Friend Request")
			}
		},
	)
	return container.NewPadded(usersList)
}

//	func showPopup(win fyne.Window) {
//		dialog.ShowInformation("Notification", "Friend request has been accepted", win)
//	}

var (
	hubsList    *widget.List
	currentHubs []api.Hub
)

func createHubsUI() fyne.CanvasObject {
	hubsList = widget.NewList(
		func() int { return len(currentHubs) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Hub"),
				layout.NewSpacer(),
				widget.NewButton("Open", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// hub := currentHubs[i]
			container := o.(*fyne.Container)

			nameLabel := container.Objects[0].(*widget.Label)
			nameLabel.SetText(currentHubs[i].HubName)

			button := container.Objects[2].(*widget.Button)
			button.OnTapped = func() {
				log.Println("Open hub")
				createHubDialog(currentHubs[i], myWindow)
			}
		},
	)
	return container.NewBorder(nil, nil, nil, nil, hubsList)
}

func createHubDialog(hub api.Hub, myWindow fyne.Window) {
	messages := widget.NewList(
		func() int { return len(hub.Messages) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Message"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			message := hub.Messages[i]
			container := o.(*fyne.Container)

			messageLabel := container.Objects[0].(*widget.Label)
			messageLabel.SetText(message.Contents)
		},
	)

	entry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Entry", Widget: entry}},
	}
	form.Refresh()

	dialog := dialog.NewCustomWithoutButtons(hub.HubName, container.NewBorder(nil, form, nil, nil, messages), myWindow)
	close := widget.NewButton("Close", func() {
		dialog.Hide()
	})
	send := widget.NewButton("Send", func() {
		if strings.TrimSpace(entry.Text) != "" {

			log.Println("Form submitted:", entry.Text)
			middle.SendHubMessage(hub.HubID, entry.Text)
			hub.Messages = append(hub.Messages, api.Message{Author: personal.GetSelf().User.UserID, Contents: entry.Text})
			messages.Refresh()
			entry.SetText("")
		}
	})
	dialog.SetButtons([]fyne.CanvasObject{container.NewBorder(nil, nil, close, send)})
	dialog.Resize(fyne.NewSize(500, 350))
	dialog.Show()
}

func onRefreshHubs(hubs []api.Hub) {
	if hubsList != nil {
		currentHubs = hubs
		hubsList.Refresh()
	}
}

var myWindow fyne.Window

func InitLoginForm(callback func(name, interest, contact string, profileImageReader io.ReadCloser)) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("What is your name?")

	interestsEntry := widget.NewEntry()
	interestsEntry.SetPlaceHolder("What interests would you like to share?")

	extraMessageEntry := widget.NewEntry()
	extraMessageEntry.SetPlaceHolder("For friends to reach you. [optional]")

	var profileImageReader io.ReadCloser
	profileImageButton := widget.NewButton("Add profile image", func() {
		profileImageDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Printf("Failed to select profile image: %v", err)
			} else {
				profileImageReader = reader
			}
		}, myWindow)

		profileImageDialog.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))

		profileImageDialog.Show()
	})

	loginForm := dialog.NewForm(
		"Login",
		"Submit",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Interests", interestsEntry),
			widget.NewFormItem("Profile Image", profileImageButton),
			widget.NewFormItem("Contact", extraMessageEntry),
		},
		func(ok bool) {
			if ok {
				callback(nameEntry.Text, interestsEntry.Text, extraMessageEntry.Text, profileImageReader)
			} else {
				myWindow.Close()
			}
		},
		myWindow,
	)

	loginForm.Resize(fyne.NewSize(500, 400))
	loginForm.Show()
}

func Init() {
	middle.Pass(onRefreshFriends, onRefreshUsers, onRefreshHubs)

	regularFont := resourceInter24ptBoldTtf

	myApp := app.New()
	myWindow = myApp.NewWindow("Agent of Friends")
	myWindow.Resize(fyne.NewSize(800, 600))
	myApp.Settings().SetTheme(&customTheme{
		Theme:           theme.DefaultTheme(),
		regularFont:     regularFont,
		primaryColor:    primaryColor,
		backgroundColor: backgroundColor,
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Peers", createUsersUI()),
		container.NewTabItem("Friends", createFriendsUI()),
		container.NewTabItem("Hubs", createHubsUI()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
}

func Run() {
	
	myWindow.ShowAndRun()
}
