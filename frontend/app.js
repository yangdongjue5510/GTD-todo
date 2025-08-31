document.addEventListener("DOMContentLoaded", () => {
    const inputField = document.querySelector(".input-row input");
    const descriptionField = document.querySelector("#description"); // New description input field
    const addButton = document.querySelector(".input-row button");
    const thingList = document.querySelector(".thing-list");

    // Fetch and display existing things from the backend
    async function fetchThings() {
        try {
            const response = await fetch("http://localhost:8080/things/");
            if (!response.ok) throw new Error("Failed to fetch things");
            const things = await response.json();
            thingList.innerHTML = ""; // Clear the list
            things.forEach((thing) => {
                const newItem = document.createElement("div");
                newItem.className = "thing-item";
                newItem.innerHTML = `<strong>${thing.title}</strong><br><em>${thing.description || "No description"}</em><br>Status: ${thing.status}`;
                thingList.appendChild(newItem);
            });
        } catch (error) {
            console.error("Error fetching things:", error);
        }
    }

    // Add new item to the backend and update the list
    async function addThing(title, description) {
        try {
            const response = await fetch("http://localhost:8080/things/", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ title, description, status: 0 }),
            });
            if (!response.ok) throw new Error("Failed to add thing");
            await fetchThings(); // Refresh the list
        } catch (error) {
            console.error("Error adding thing:", error);
        }
    }

    // Event listener for adding a new thing
    addButton.addEventListener("click", () => {
        const inputValue = inputField.value.trim();
        const descriptionValue = descriptionField.value.trim();
        if (inputValue) {
            addThing(inputValue, descriptionValue);
            inputField.value = "";
            descriptionField.value = "";
        }
    });

    // Allow pressing Enter to add item
    inputField.addEventListener("keypress", (event) => {
        if (event.key === "Enter") {
            addButton.click();
        }
    });

    // Initial fetch of things
    fetchThings();
});

// Updated reference to the correct HTML file
// Ensure the frontend is linked to `index.html` instead of `html.html`
