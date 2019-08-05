function openPopUp(errorType) {
    $("#wrap").css("filter", "brightness(50%)");
    var error_string
    switch (errorType) {
        case "empty":
            error_string = "<p>Введите логин или пароль</p>"
            break;
        case "badLogin":
            error_string = "<p>Неправильно введен пароль или логин</p>"
            break;
        default:
            error_string = "<p>Интересная и неопознаная ошибка</p>"
            break;
    }
    document.getElementById("error_text").innerHTML = error_string;
    $(".js-overlay").fadeIn();
};

$(".js-close-popup").click(function() {
    $(".js-overlay").fadeOut();
    $("#wrap").css("filter", "none");
});

$(document).mouseup(function(e) {
    var popup = $(".js-popup");
    if (e.target != popup[0] && popup.has(e.target).length === 0) {
        $(".js-overlay").fadeOut();
        $("#wrap").css("filter", "none");
    }
});

function auth(form) {
    var username = form.username.value;
    var password = form.password.value;

    if (username.length == 0 || password.length == 0) {
        openPopUp("empty");
    } else {
        $.ajax({
            url: "/login",
            method: "POST",
            data: {
                username: username,
                password: password,
            },
            success: function() {
                location.href = "/profile/" + username;
            },
            error: function() {
                openPopUp("badLogin");
                form.reset();
            }
        })
    }
};