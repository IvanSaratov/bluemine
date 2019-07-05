function tClose(){
    var href = location.href.split('/');
    var id = href[5];
    if (id.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/tasks/close",
            method: "POST",
            data: {
                id: id
            },
            success: function(){
                location.href = "/tasks";
            }
        });
    }
};