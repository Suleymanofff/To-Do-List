document.addEventListener('DOMContentLoaded', loadTasks)

async function loadTasks() {
	const res = await fetch('/tasks')
	const tasks = await res.json()

	const taskList = document.getElementById('taskList')
	taskList.innerHTML = ''

	tasks.forEach(task => {
		const li = document.createElement('li')
		li.innerHTML = `
            ${task.text}
            <button onclick="deleteTask(${task.id})">Delete</button>
        `
		taskList.appendChild(li)
	})
}

async function addTask() {
	const input = document.getElementById('taskInput')
	const text = input.value.trim()

	if (!text) return

	await fetch('/add', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ text }),
	})

	input.value = ''
	loadTasks()
}

async function deleteTask(id) {
	await fetch('/delete', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ id }),
	})

	loadTasks()
}
