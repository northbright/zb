
var name = "";
var tel = "";
var grade = "";
var campus = "";
var availabeCampuses = [];
var currentPeriod = "";
var wantedPeriod1 = "";
var wantedPeriod2 = ""; 

function getGrades() {
  $('#grade').html('');

  $.ajax({
        url: '/grades',
        type: 'GET',
        dataType: 'json',
        error : function () {
            alert("获取年级信息失败");
        },   
        success: function (data) {
            if (!data.success) {
                alert("获取年级信息失败:" + data.err_msg);
            } else {
                $.each(data.grades, function(index, value){
                    $('#grade').append('<option value="' + value + '">' + value + '</option>');
                });    
                $("#grade").val(data.grades[0]).change();
            }
        }
  });
}

function getCampuses(grade) {
  $('#campus').html('');

  $.ajax({
        url: '/campus/' + grade,
        type: 'GET',
        dataType: 'json',
        error : function () {
            alert("获取校区信息失败");
        },
        success: function (data) {
            if (!data.success) {
                alert("获取校区信息失败:" + data.err_msg);
            } else {
                $.each(data.campuses, function(index, value){
                    $('#campus').append('<option value="' + value + '">' + value + '</option>');
                });
                $("#campus").val(data.campuses[0]).change();
            }
        }
  });
}

$(document).ready(function () {
  //alert("ready");
});

$(document).on("pageinit","#page1",function(){
  //alert("page 1");
  // Clear all options in select grade.
  $('#grade').html('');

  $('#grade').change(function () {
    grade = $(this).find('option:selected').val();
    //alert(grade);
    getCampuses(grade);
  });

  getGrades();
});


