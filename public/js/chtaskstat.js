function taskReOpen(){
    var id = $('.pagetitle').attr("id");

    if (id.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/task/open",
            method: "POST",
            data: {
                id: id
            },
            success: function(){
                location.href = "/task/" + id;
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
            url: "/task/close",
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