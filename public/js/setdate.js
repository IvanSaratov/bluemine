function setDate() {
    var today = new Date();

    document.getElementById("input_date_start").value = today.getFullYear() + '-' + ('0' + (today.getMonth() + 1)).slice(-2) + '-' + ('0' + today.getDate()).slice(-2);
}