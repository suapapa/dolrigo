var width = 600;
var height = 600;
var radius = Math.min(width, height) / 2;
var color_table = d3.scale.category20();

var gradient_colors = [
	["#ff6600", "#cc6600"],
	["#cc6600", "#996600"],
	["#996600", "#666600"],
	["#666600", "#669933"],
	["#669933", "#336600"],
	["#336600", "#336633"],
	["#336633", "#009966"],
	["#009966", "#009999"],
	["#009999", "#006699"],
	["#006699", "#003366"],
	["#003366", "#003399"],
	["#003399", "#330099"],
	["#330099", "#330066"],
	["#330066", "#660099"],
	["#660099", "#990066"],
	["#990066", "#660033"],
	["#660033", "#cc3333"],
	["#cc3333", "#993300"],
	["#993300", "#cc3300"],
	["#cc3300", "#ff6600"],
	["#df5649", "#c10e00"],
	["#e08a3e", "#c94c00"],
	["#c4ac3b", "#987600"],
	["#99ac50", "#597400"],
	["#62af82", "#1f7f41"],
	["#5ca3c0", "#0e6894"],
	["#808bcf", "#3e4aa9"],
	["#ab74d2", "#782caa"],
	["#ca6a8b", "#a1244f"],
	["#a5a5a5", "#6c6b6b"],
	["#f55a89", "#c83461"]
];
var svg;
var wheel;
var top_degree;

var friction_ratio = 0.03; // 50 ms 당 wheel_speed 감소 비율
var friction_min = 0.05; // 50 ms 당 wheel_speed 감소 최소값

function degree2idx(items, degree) {
	var divisor = 360 / items.length;

	return Math.floor((degree + 180 / items.length) / divisor) % items.length;
}

function toggle_photo(items) {
	var focused_idx = degree2idx(items, top_degree);
	var focused_item = $("#photo_container .focused");

	if (items[focused_idx] != focused_item.data())
	{
		focused_item.removeClass("focused");
		$("#photo_container .photo:eq(" + focused_idx + ")").addClass("focused");
	}
}

function generate_arc_path(max_items) {
	var arc_path_generator = 
		d3.svg.arc()
			.outerRadius(radius - 10)
			.innerRadius(radius / 2)
			.startAngle(-Math.PI/max_items)
			.endAngle(Math.PI/max_items);

	return arc_path_generator(null);
}

function prepare_wheel(items) {
	top_degree = 0;

	svg = d3.select("#wheel_container").append("svg")
		.attr("width", 600)
		.attr("height", 600);

	var gradient = svg.selectAll(".gradient")
			.data(gradient_colors)
		.enter().append("svg:defs")		
			.append("svg:linearGradient")
			.attr("id", function(d, i) { return "grad"+i; })
			.attr("x1", "0%")
			.attr("y1", "100%")
			.attr("x2", "0%")
			.attr("y2", "0%")
			.attr("spreadMethod", "pad")
			.selectAll(".stop")
					.data(function(d) { return d; })
				.enter().append("svg:stop")
					.attr("offset", function(d, i) { if (i == 0) return "0%"; else return "100%"; })
					.attr("stop-color", function(d) { return d; })
					.attr("stop-opacity", 1);

	svg.append("svg:defs")		
			.append("svg:radialGradient")
			.attr("id", "grad_center")
			.attr("x1", "0%")
			.attr("y1", "100%")
			.attr("x2", "0%")
			.attr("y2", "0%")
			.attr("spreadMethod", "pad")
			.selectAll(".stop")
				.data(["#deb887", "#b8860b"])
				.enter().append("svg:stop")
					.attr("offset", function(d, i) { if (i == 0) return "0%"; else return "100%"; })
					.attr("stop-color", function(d) { return d; })
					.attr("stop-opacity", 1);

	wheel = svg.append("g")
		.attr("transform", "translate(" + width / 2 + "," + height / 2 + ") rotate(" + top_degree +")");

	var img = 
		wheel.append("g")
			.attr("class", "img");
	img.append("image")
		.attr("xlink:href", "img/bg.png")
		.attr("x", -100)
		.attr("y", -100)
		.attr("width", 200)
		.attr("height", 200);

/*
	var title = 
		wheel.append("g")
			.attr("class", "title");

	title.append("circle")
		.attr("r", radius / 2 - 5)
		.style("fill", "url('#grad_center')");

	title.append("text")
		.attr("y", -50)
		.attr("dy", ".35em")
		.style("text-anchor", "middle")
		.text("(환)Egslee님(영)");

	title.append("text")
		.attr("dy", ".35em")
		.style("text-anchor", "middle")
		.text("개발자 노가리방");

	title.append("text")
		.attr("dy", ".35em")
		.attr("y", 50)
		.style("text-anchor", "middle")
		.text("2022 워크샵");

*/ 

	var arc_path = generate_arc_path(items.length);

	var g = 
		wheel.selectAll(".arc")
				.data(items)
			.enter().append("g")
				.attr("class", "arc")
				.attr("transform", 
					function (d, i) { 
						return "rotate(" + (360 * i / items.length) + ")";
					});

	g.append("path")
			.attr("d", arc_path)
			.style("fill", 
				function(d, i) { 
					return "url('#grad" + (i % gradient_colors.length) + "')";
					//return color_table(i); 
				}
			);

	g.append("text")
			.attr("transform", 
				function(d, i) { 
					return "translate(0, -" + (radius - 10) + ")";
				})
			.attr("dy", ".35em")
			.text(function(d) { return d.name; });
}

function spin_wheel(items, wheel_speed) {

	if (wheel_speed <= 0) {
		setTimeout(function () { display_focused_person(items); }, 50);
		return;
	}

	top_degree += wheel_speed;

	wheel.transition()
		.duration(50)
		.attr("transform", "translate(" + width / 2 + "," + height / 2 + ") rotate(" + (- top_degree) +")")
		.each("end", 
			function () { 
				toggle_photo(items);
				wheel_speed -= Math.max(wheel_speed * friction_ratio, friction_min);
				spin_wheel(items, wheel_speed); 
			});
}

function display_focused_person(items) {
	top_degree -= Math.floor(top_degree / 360) * 360;

	toggle_photo(items);

	var item = $("#photo_container .focused").data();

	alert("당첨자: " + item.name);
}

