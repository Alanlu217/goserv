document.querySelectorAll('.post-button').forEach(button => {
    button.addEventListener('click', async function() {
        const button_id = this.getAttribute('button_id');
        
        const response = await fetch('/dash', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                id: Number(button_id)
            })
        });

        console.log(response.body)

        if (!response.ok) {
            alert("Request Failed")
        }
    });
});

