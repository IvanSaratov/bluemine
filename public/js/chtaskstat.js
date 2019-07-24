function taskReOpen(){
    var id = $('.pagetitle').attr("id");

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

function taskClose(){
    var id = $('.pagetitle').attr("id");

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