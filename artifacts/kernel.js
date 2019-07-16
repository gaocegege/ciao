define(function(){
    var nb = IPython.notebook;
    window.CodeCell = require("notebook/js/codecell").CodeCell;

    return {onload: function(){
      console.info('load cell');

      events.on("create.Cell", function(evt, param) {
        add_lan_selector(param.cell);
      });
  
    }}

    function add_lan_selector(cell) {
        // add a new one
        var select = $("<select/>")
        .attr("id", "cell_kernel_selector")
        .css("right", "100px")
        .attr("class", "select-xs cell_kernel_selector");
        select.append(
        $("<option/>")
            .attr("value", "1")
            .text("1 PS, 2 Worker")
        );
        select.append(
            $("<option/>")
                .attr("value", "1")
                .text("1 PS, 4 Worker")
            );
        cell.element.find("div.input_area").prepend(select);
        return select;
    }
  
  });
