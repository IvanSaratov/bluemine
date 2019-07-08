$(document).ready(function() {
    $('#iconew').hover(
        function() {
            $('#new_item').css('display', 'block');
        },
        function() {
            $('#new_item').css('display', 'none');
        }
    )
});