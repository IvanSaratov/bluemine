function gChange() {
    var href = location.href.split('/');
    var id = href[5];
    $.get("/groupchange", { id: id }).done(function(data){
        $('#new_group').show(300);
        $('#input_group_name').val(data.GroupName)

        for (i = 0; i < data.GroupMembers.length; i++) {
            $("input[value='"+data.GroupMembers[i].UserName+"']").prop('checked', true);
        };
    })
};
