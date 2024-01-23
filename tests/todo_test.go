package tests

//Introduction to testing.  Note that testing is built into go and we will be using
//it extensively in this class. Below is a starter for your testing code.  In
//addition to what is built into go, we will be using a few third party packages
//that improve the testing experience.  The first is testify.  This package brings
//asserts to the table, that is much better than directly interacting with the
//testing.T object.  Second is gofakeit.  This package provides a significant number
//of helper functions to generate random data to make testing easier.

import (
	"fmt"
	"os"
	"testing"

	"drexel.edu/todo/db"
	fake "github.com/brianvoe/gofakeit/v6" //aliasing package name
	"github.com/stretchr/testify/assert"
)

// Note the default file path is relative to the test package location.  The
// project has a /tests path where you are at and a /data path where the
// database file sits.  So to get there we need to back up a directory and
// then go into the /data directory.  Thus this is why we are setting the
// default file name to "../data/todo.json"
const (
	DEFAULT_DB_FILE_NAME = "../data/todo.json"
)

var (
	DB *db.ToDo
)

// note init() is a helpful function in golang.  If it exists in a package
// such as we are doing here with the testing package, it will be called
// exactly once.  This is a great place to do setup work for your tests.
func init() {
	//Below we are setting up the gloabal DB variable that we can use in
	//all of our testing functions to make life easier
	testdb, err := db.New(DEFAULT_DB_FILE_NAME)
	if err != nil {
		fmt.Print("ERROR CREATING DB:", err)
		os.Exit(1)
	}

	DB = testdb //setup the global DB variable to support test cases

	//Now lets start with a fresh DB with the sample test data
	testdb.RestoreDB()
}

// Sample Test, will always pass, comparing the second parameter to true, which
// is hard coded as true
func TestTrue(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestAddHardCodedItem(t *testing.T) {
	item := db.ToDoItem{
		Id:     999,
		Title:  "This is a test case item",
		IsDone: false,
	}
	t.Log("Testing Adding a Hard Coded Item: ", item)

	//TODO: finish this test, add an item to the database and then
	//check that it was added correctly by looking it back up
	//use assert.NoError() to ensure errors are not returned.
	//explore other useful asserts in the testify package, see
	//https://github.com/stretchr/testify.  Specifically look
	//at things like assert.Equal() and assert.Condition()

	//I will get you started, uncomment the lines below to add to the DB
	//and ensure no errors:
	//---------------------------------------------------------------
	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to DB")

	//TODO: Now finish the test case by looking up the item in the DB
	//and making sure it matches the item that you put in the DB above

	toDoItem, getItemErr := DB.GetItem(item.Id)
	assert.NoError(t, getItemErr, "Error getting item from DB")
	assert.Equal(t, item, toDoItem, "Item added to DB does not match item retrieved from DB")

	deleteItemErr := DB.DeleteItem(item.Id)
	assert.NoError(t, deleteItemErr, "Error deleting item from DB")
}

func TestAddRandomStructItem(t *testing.T) {
	//You can also use the Stuct() fake function to create a random struct
	//Not going to do anyting
	item := db.ToDoItem{}
	err := fake.Struct(&item)
	t.Log("Testing Adding a Randomly Generated Struct: ", item)
	assert.NoError(t, err, "Created fake item OK")

	//TODO: Complete the test
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	toDoItem, getItemErr := DB.GetItem(item.Id)
	assert.NoError(t, getItemErr, "Error getting item from DB")
	assert.Equal(t, item, toDoItem, "Item added to DB does not match item retrieved from DB")

	deleteItemErr := DB.DeleteItem(item.Id)
	assert.NoError(t, deleteItemErr, "Error deleting item from DB")
}

func TestAddRandomItem(t *testing.T) {
	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	t.Log("Testing Adding an Item with Random Fields: ", item)
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	toDoItem, getItemErr := DB.GetItem(item.Id)
	assert.NoError(t, getItemErr, "Error getting item from DB")
	assert.Equal(t, item, toDoItem, "Item added to DB does not match item retrieved from DB")

	deleteItemErr := DB.DeleteItem(item.Id)
	assert.NoError(t, deleteItemErr, "Error deleting item from DB")
}

//TODO: Create additional tests to showcase the correct operation of your program
//for example getting an item, getting all items, updating items, and so on. Be
//creative here.

func TestRestoreDB(t *testing.T) {

	//GetAllItems from DB
	originalTodoList, getAllItemErr := DB.GetAllItems()
	assert.NoError(t, getAllItemErr, "Error getting all items from DB")

	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	// Add item to DB
	t.Log("Adding new todo item: ", item)
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	//Now lets try to get the item and validate if item was added correctly
	newToDoItem, newGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, newGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, newToDoItem, "Item added to DB does not match item retrieved from DB")

	//Now lets restore the DB
	restoreDBErr := DB.RestoreDB()
	assert.NoError(t, restoreDBErr, "Error restoring DB")

	//Now lets try to get the item again
	newToDoItem, newGetItemErr = DB.GetItem(item.Id)
	assert.EqualError(t, newGetItemErr, "todo trying to fetch doesnt exists", "Item should not be found in DB")
	assert.Equal(t, db.ToDoItem{}, newToDoItem, "Item should not be found in DB")

	//Now lets test that the DB was restored correctly by comparing all items
	restoredToDoList, getAllItemErr := DB.GetAllItems()
	assert.NoError(t, getAllItemErr, "Error getting all items from DB")
	assert.Equal(t, len(originalTodoList), len(restoredToDoList), "Items in DB should be the same as before updating DB")

	//Now lets test that the DB was restored correctly by comparing all items
	for _, originalToDoItem := range originalTodoList {
		originalToDoItem, originalGetItemErr := DB.GetItem(originalToDoItem.Id)
		assert.NoError(t, originalGetItemErr, "Error getting item from DB")
		restoredToDoItem, restoredGetItemError := DB.GetItem(originalToDoItem.Id)
		assert.NoError(t, restoredGetItemError, "Error getting item from DB")
		assert.Equal(t, originalToDoItem, restoredToDoItem, "Item in restored and original DB should be the same")
	}
}

func TestDeleteItem(t *testing.T) {

	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	// Add item to DB
	t.Log("Adding new todo item: ", item)
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	//Now lets try to get the item and validate if item was added correctly
	newToDoItem, newGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, newGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, newToDoItem, "Item added to DB does not match item retrieved from DB")

	//Now lets delete the item
	deleteItemErr := DB.DeleteItem(item.Id)
	assert.NoError(t, deleteItemErr, "Error deleting item from DB")

	//Now lets try to get the item again
	newToDoItem, newGetItemErr = DB.GetItem(item.Id)
	assert.EqualError(t, newGetItemErr, "todo trying to fetch doesnt exists", "Item should not be found in DB")
	assert.Equal(t, db.ToDoItem{}, newToDoItem, "Item should not be found in DB")

	//Now lets try to delete the item again
	deleteItemErr = DB.DeleteItem(item.Id)
	assert.EqualError(t, deleteItemErr, "todo trying to delete doesnt exists", "Item should not be found in DB")
}

func TestUpdateItem(t *testing.T) {

	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	// Add item to DB
	t.Log("Adding new todo item: ", item)
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	//Now lets try to get the item and validate if item was added correctly
	newToDoItem, newGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, newGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, newToDoItem, "Item added to DB does not match item retrieved from DB")

	//Lets update the item
	item.Title = fake.JobTitle()
	item.IsDone = fake.Bool()

	//Now lets update the item to the DB
	updateItemErr := DB.UpdateItem(item)
	assert.NoError(t, updateItemErr, "Error updating item in DB")

	//Now lets try to get the item and validate if item was added correctly
	updatedToDoItem, updatedGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, updatedGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, updatedToDoItem, "Item added to DB does not match item retrieved from DB")

	//Now lets Restore the DB to original state so that other tests can run
	restoreDBErr := DB.RestoreDB()
	assert.NoError(t, restoreDBErr, "Error restoring DB")
}

func TestGetItem(t *testing.T) {

	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	// Add item to DB
	t.Log("Adding new todo item: ", item)
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	//Now lets try to get the item and validate if item was added correctly
	newToDoItem, newGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, newGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, newToDoItem, "Item added to DB does not match item retrieved from DB")

	//Now lets Restore the DB to original state so that other tests can run
	restoreDBErr := DB.RestoreDB()
	assert.NoError(t, restoreDBErr, "Error restoring DB")

	//Now lets try to get the item again
	newToDoItem, newGetItemErr = DB.GetItem(item.Id)
	assert.EqualError(t, newGetItemErr, "todo trying to fetch doesnt exists", "Item should not be found in DB")
	assert.Equal(t, db.ToDoItem{}, newToDoItem, "Item should not be found in DB")
}

func TestGetAllItems(t *testing.T) {

	//GetAllItems from DB
	originalTodoList, getAllItemErr := DB.GetAllItems()
	assert.NoError(t, getAllItemErr, "Error getting all items from DB")

	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	// Add item to DB
	t.Log("Adding new todo item: ", item)
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	//Now lets try to get the item and validate if item was added correctly
	newToDoItem, newGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, newGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, newToDoItem, "Item added to DB does not match item retrieved from DB")

	//Now lets get all items from DB
	newTodoList, getAllItemErr := DB.GetAllItems()
	assert.NoError(t, getAllItemErr, "Error getting all items from DB")

	//Now lets test that the DB was restored correctly by comparing length of list
	assert.Equal(t, len(originalTodoList)+1, len(newTodoList), "Length of list should be one more than original list")

	//Now lets Restore the DB to original state so that other tests can run
	restoreDBErr := DB.RestoreDB()
	assert.NoError(t, restoreDBErr, "Error restoring DB")
}

func TestChangeItemDoneStatus(t *testing.T) {

	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	// Add item to DB
	t.Log("Adding new todo item: ", item)
	addItemErr := DB.AddItem(item)
	assert.NoError(t, addItemErr, "Error adding item to DB")

	//Now lets try to get the item and validate if item was added correctly
	newToDoItem, newGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, newGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, newToDoItem, "Item added to DB does not match item retrieved from DB")

	//Now lets change the done status of the item
	item.IsDone = !item.IsDone

	//Now lets update the item to the DB
	updateItemErr := DB.UpdateItem(item)
	assert.NoError(t, updateItemErr, "Error updating item in DB")

	//Now lets try to get the item and validate if item was added correctly
	updatedToDoItem, updatedGetItemErr := DB.GetItem(item.Id)
	assert.NoError(t, updatedGetItemErr, "Error getting item from DB")
	assert.Equal(t, item, updatedToDoItem, "Item added to DB does not match item retrieved from DB")

	//Now lets Restore the DB to original state so that other tests can run
	restoreDBErr := DB.RestoreDB()
	assert.NoError(t, restoreDBErr, "Error restoring DB")
}
