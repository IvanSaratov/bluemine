function tOpen(){
    var href = location.href.split('/');
    var id = href[5];
    if (id.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/tasks/open",
            method: "POST",
            data: {
                id: id
            },
            success: function(){
                location.href = "/tasks/show/" + id;
            }
        });
    }
};