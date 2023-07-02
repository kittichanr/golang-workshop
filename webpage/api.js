function getSampleData() {
    // Make an API request
    fetch('https://e97e-2405-9800-bc10-509e-8c5f-be4a-8547-f190.ap.ngrok.io/course')
        .then(response => response.json())
        .then(data => {
            appendData(data)
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function appendData(data) {
    const dataDiv = document.getElementById('data');
    for (let i = 0; i < data.length; i++) {
        const div = document.createElement('div')
        div.innerHTML = "CourseID: " + data[i].id + ' ' + data[i].name + ' ' + data[i].price + ' ' + data[i].instructor
        dataDiv.appendChild(div)
    }
}