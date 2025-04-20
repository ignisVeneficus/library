// https://tailwindcolor.com/ 800,600,400
const colors =[
    //800
"#991B1B",
"#9A3412",
"#92400E",
"#854D0E",
"#3F6212",
"#166534",
"#065F46",
"#115E59",
"#155E75",
"#075985",
"#1E40AF",
"#3730A3",
"#5B21B6",
"#6B21A8",
"#86198F",
"#9D174D",
"#9F1239",
"#292524",
"#1E293B",
    //600
"#DC2626",
"#EA580C",
"#D97706",
"#CA8A04",
"#65A30D",
"#16A34A",
"#059669",
"#0D9488",
"#0891B2",
"#0284C7",
"#2563EB",
"#4F46E5",
"#7C3AED",
"#9333EA",
"#C026D3",
"#DB2777",
"#E11D48",
"#57534E",
"#475569"

];


//https://stackoverflow.com/questions/5623838/rgb-to-hex-and-hex-to-rgb
function hexToRgb(hex) {
    var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? [
        parseInt(result[1], 16),
        parseInt(result[2], 16),
        parseInt(result[3], 16)
    ] : null;
}

//https://stackoverflow.com/questions/11867545/change-text-color-based-on-brightness-of-the-covered-background-area
function setContrast(rgbColor) {
    // http://www.w3.org/TR/AERT#color-contrast
    const brightness = Math.round(((parseInt(rgbColor[0]) * 299) +
        (parseInt(rgbColor[1]) * 587) +
        (parseInt(rgbColor[2]) * 114)) / 1000);
    return (brightness > 125) ? 'black' : 'white';
}