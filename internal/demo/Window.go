package demo

import (
	"fmt"
	"strings"

	"github.com/AllenDang/cimgui-go"
)

type windowFlags struct {
	noTitlebar     bool
	noScrollbar    bool
	noMenu         bool
	noMove         bool
	noResize       bool
	noCollapse     bool
	noNav          bool
	noBackground   bool
	noBringToFront bool
}

func (f windowFlags) combined() cimgui.ImGuiWindowFlags {
	var flags cimgui.ImGuiWindowFlags = cimgui.ImGuiWindowFlags_None
	if f.noTitlebar {
		flags |= cimgui.ImGuiWindowFlags_NoTitleBar
	}
	if f.noScrollbar {
		flags |= cimgui.ImGuiWindowFlags_NoScrollbar
	}
	if !f.noMenu {
		flags |= cimgui.ImGuiWindowFlags_MenuBar
	}
	if f.noMove {
		flags |= cimgui.ImGuiWindowFlags_NoMove
	}
	if f.noResize {
		flags |= cimgui.ImGuiWindowFlags_NoResize
	}
	if f.noCollapse {
		flags |= cimgui.ImGuiWindowFlags_NoCollapse
	}
	if f.noNav {
		flags |= cimgui.ImGuiWindowFlags_NoNav
	}
	if f.noBackground {
		flags |= cimgui.ImGuiWindowFlags_NoBackground
	}
	if f.noBringToFront {
		flags |= cimgui.ImGuiWindowFlags_NoBringToFrontOnFocus
	}
	return flags
}

var window = struct {
	flags   windowFlags
	noClose bool

	widgets widgets
	layout  layout
	popups  popups
	columns columns
	tables  tables
	misc    misc
}{}

func bulletText(text string) {
	cimgui.Bullet()
	cimgui.Text(text)
}

// Show demonstrates most ImGui features that were ported to Go.
// This function tries to recreate the original demo window as closely as possible.
//
// In theory, if both windows would provide the identical functionality, then the wrapper would be complete.
func Show(keepOpen *bool) {

	cimgui.SetNextWindowPos(cimgui.ImVec2{X: 650, Y: 20}, cimgui.ImGuiCond_FirstUseEver, cimgui.ImVec2{})
	cimgui.SetNextWindowSize(cimgui.ImVec2{X: 550, Y: 680}, cimgui.ImGuiCond_FirstUseEver)

	if window.noClose {
		keepOpen = nil
	}
	if !cimgui.Begin("ImGui-Go Demo", keepOpen, window.flags.combined()) {
		// Early out if the window is collapsed, as an optimization.
		cimgui.End()
		return
	}

	// Use fixed width for labels (by passing a negative value), the rest goes to widgets.
	// We choose a width proportional to our font size.
	cimgui.PushItemWidth(cimgui.GetFontSize() * -12)

	// MenuBar
	if cimgui.BeginMenuBar() {
		if cimgui.BeginMenu("Menu", true) {
			cimgui.EndMenu()
		}
		if cimgui.BeginMenu("Examples", true) {
			cimgui.EndMenu()
		}
		if cimgui.BeginMenu("Tools", true) {
			cimgui.EndMenu()
		}

		cimgui.EndMenuBar()
	}

	cimgui.Text(fmt.Sprintf("dear cimgui says hello. (%s)", cimgui.GetVersion()))
	cimgui.Spacing()

	if cimgui.CollapsingHeader_TreeNodeFlags("Help", cimgui.ImGuiTreeNodeFlags_None) {
		cimgui.Text("ABOUT THIS DEMO:")
		bulletText("Sections below are demonstrating many aspects of the wrapper.")
		bulletText("This demo may not be complete. Refer to the \"native\" demo window for a full overview.")
		bulletText("The \"Examples\" menu above leads to more demo contents.")
		bulletText("The \"Tools\" menu above gives access to: About Box, Style Editor,\n" +
			"and Metrics (general purpose Dear ImGui debugging tool).")
		cimgui.Separator()

		cimgui.Text("PROGRAMMER GUIDE:")
		bulletText("See the demo.Show() code in internal/demo/Window.go. <- you are here!")
		bulletText("See comments in cimgui.cpp.")
		bulletText("See example applications in the examples/ folder.")
		bulletText("Read the FAQ at http://www.dearimgui.org/faq/")
		bulletText("Set 'io.ConfigFlags |= NavEnableKeyboard' for keyboard controls.")
		bulletText("Set 'io.ConfigFlags |= NavEnableGamepad' for gamepad controls.")
		cimgui.Separator()

		cimgui.Text("USER GUIDE:")
		showUserGuide()
	}

	// MISSING: Configuration

	if cimgui.CollapsingHeader_TreeNodeFlags("Window options", cimgui.ImGuiTreeNodeFlags_None) {
		cimgui.Checkbox("No titlebar", &window.flags.noTitlebar)
		cimgui.SameLine(150, -1)
		cimgui.Checkbox("No scrollbar", &window.flags.noScrollbar)
		cimgui.SameLine(300, -1)
		cimgui.Checkbox("No menu", &window.flags.noMenu)
		cimgui.Checkbox("No move", &window.flags.noMove)
		cimgui.SameLine(150, -1)
		cimgui.Checkbox("No resize", &window.flags.noResize)
		cimgui.SameLine(300, -1)
		cimgui.Checkbox("No collapse", &window.flags.noCollapse)
		cimgui.Checkbox("No close", &window.noClose)
		cimgui.SameLine(150, -1)
		cimgui.Checkbox("No nav", &window.flags.noNav)
		cimgui.SameLine(300, -1)
		cimgui.Checkbox("No background", &window.flags.noBackground)
		cimgui.Checkbox("No bring to front", &window.flags.noBringToFront)
	}

	// All demo contents
	window.widgets.show()
	window.layout.show()
	window.popups.show()
	window.columns.show()
	window.tables.show()
	window.misc.show()

	// End of ShowDemoWindow()
	cimgui.End()
}

func showUserGuide() {
	bulletText("Double-click on title bar to collapse window.")
	bulletText("Click and drag on lower corner to resize window\n(double-click to auto fit window to its contents).")
	bulletText("CTRL+Click on a slider or drag box to input value as text.")
	bulletText("TAB/SHIFT+TAB to cycle through keyboard editable fields.")

	// MISSING: Allow FontUserScaling

	bulletText("While inputing text:\n")
	cimgui.Indent(0)
	bulletText("CTRL+Left/Right to word jump.")
	bulletText("CTRL+A or double-click to select all.")
	bulletText("CTRL+X/C/V to use clipboard cut/copy/paste.")
	bulletText("CTRL+Z,CTRL+Y to undo/redo.")
	bulletText("ESCAPE to revert.")
	bulletText("You can apply arithmetic operators +,*,/ on numerical values.\nUse +- to subtract.")
	cimgui.Unindent(0)
	bulletText("With keyboard navigation enabled:")
	cimgui.Indent(0)
	bulletText("Arrow keys to navigate.")
	bulletText("Space to activate a widget.")
	bulletText("Return to input text into a widget.")
	bulletText("Escape to deactivate a widget, close popup, exit child window.")
	bulletText("Alt to jump to the menu layer of a window.")
	bulletText("CTRL+Tab to select a window.")
	cimgui.Unindent(0)
}

type widgets struct {
	buttonClicked int
	check         bool
	radio         int32
}

// nolint: nestif
func (widgets *widgets) show() {
	if !cimgui.CollapsingHeader_TreeNodeFlags("Widgets", cimgui.ImGuiTreeNodeFlags_None) {
		return
	}

	if cimgui.TreeNode_Str("Basic") {
		if cimgui.Button("Button", cimgui.ImVec2{}) {
			widgets.buttonClicked++
		}
		if widgets.buttonClicked&1 != 0 {
			cimgui.SameLine(0, -1)
			cimgui.Text("Thanks for clicking me!")
		}

		cimgui.Checkbox("checkbox", &widgets.check)

		if cimgui.RadioButton_Bool("radio a", widgets.radio == 0) {
			widgets.radio = 0
		}
		cimgui.SameLine(0, -1)
		if cimgui.RadioButton_Bool("radio b", widgets.radio == 1) {
			widgets.radio = 1
		}
		cimgui.SameLine(0, -1)
		if cimgui.RadioButton_Bool("radio c", widgets.radio == 2) {
			widgets.radio = 2
		}

		str := strings.Join([]string{
			"one item",
			"two items",
			"three items",
		}, "\x00") + "\x00"
		cimgui.Combo_Str("combo", &widgets.radio, str, -1)

		cimgui.TreePop()
	}
}

type tables struct {
	background     bool
	borders        bool
	noInnerBorders bool
	header         bool
}

var demoTableHeader = []string{
	"Name", "Favourite Food", "Favourite Colour",
}

var demoTable = [][]string{
	{"Eric", "Bannana", "Yellow"},
	{"Peter", "Apple", "Red"},
	{"Bruce", "Liquorice", "Black"},
	{"Aaron", "Chocolates", "Blue"},
}

func (tables *tables) show() {
	if !cimgui.CollapsingHeader_TreeNodeFlags("Tables", cimgui.ImGuiTreeNodeFlags_None) {
		return
	}

	if cimgui.TreeNode_Str("Rows & Columns") {
		if cimgui.BeginTable("tableRowsAndColumns", 3, cimgui.ImGuiTableFlags_None, cimgui.ImVec2{}, 0) {
			for row := 0; row < 4; row++ {
				cimgui.TableNextRow(cimgui.ImGuiTableRowFlags_None, 0)
				for column := int32(0); column < 3; column++ {
					cimgui.TableSetColumnIndex(column)
					cimgui.Text(fmt.Sprintf("Row %d Column %d", row, column))
				}
			}
			cimgui.EndTable()
		}
		cimgui.TreePop()
	}

	if cimgui.TreeNode_Str("Options") {
		// tables are useful for more than tabulated data. we use tables here
		// to facilitate layout of the option checkboxes
		if cimgui.BeginTable("tableOptions", 2, cimgui.ImGuiTableFlags_None, cimgui.ImVec2{}, 0) {
			cimgui.TableNextRow(cimgui.ImGuiTableRowFlags_None, 0)
			if cimgui.TableNextColumn() {
				cimgui.Checkbox("Background", &tables.background)
			}
			if cimgui.TableNextColumn() {
				cimgui.Checkbox("Header Row", &tables.header)
			}

			cimgui.TableNextRow(cimgui.ImGuiTableRowFlags_None, 0)
			if cimgui.TableNextColumn() {
				cimgui.Checkbox("Borders", &tables.borders)
			}
			if tables.borders {
				if cimgui.TableNextColumn() {
					cimgui.Checkbox("No Inner Borders", &tables.noInnerBorders)
				}
			}

			cimgui.EndTable()
		}

		// set flags according to the options that have been selected
		flgs := cimgui.ImGuiTableFlags_None
		if tables.background {
			flgs |= cimgui.ImGuiTableFlags_RowBg
		}
		if tables.borders {
			flgs |= cimgui.ImGuiTableFlags_Borders
			if tables.noInnerBorders {
				flgs |= cimgui.ImGuiTableFlags_NoBordersInBody
			}
		}

		if cimgui.BeginTable("tableRowsAndColumns", int32(len(demoTableHeader)), cimgui.ImGuiTableFlags(flgs), cimgui.ImVec2{}, 0.0) {
			if tables.header {
				cimgui.TableHeadersRow()
				for column := 0; column < len(demoTableHeader); column++ {
					cimgui.TableSetColumnIndex(int32(column))
					cimgui.Text(demoTableHeader[column])
				}
			}

			for row := 0; row < len(demoTable); row++ {
				cimgui.TableNextRow(cimgui.ImGuiTableRowFlags_None, 0)
				for column := 0; column < len(demoTableHeader); column++ {
					cimgui.TableSetColumnIndex(int32(column))
					cimgui.Text(demoTable[row][column])
				}
			}
			cimgui.EndTable()
		}
		cimgui.TreePop()
	}
}

type layout struct {
}

func (layout *layout) show() {

}

type popups struct {
}

func (popups *popups) show() {

}

type columns struct {
}

func (columns *columns) show() {

}

type misc struct {
}

func (misc *misc) show() {

}
