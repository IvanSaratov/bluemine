$(document).ready(function() {
    $('#ico_new').hover(
        function() {
            $('#new_item').css('display', 'block');
        },
        function() {
            $('#new_item').css('display', 'none');
        }
    )
    
    $('#new_item').hover(
        function() {
            $('#iconew').toggleClass("hovered", 200);
        }
    );
});