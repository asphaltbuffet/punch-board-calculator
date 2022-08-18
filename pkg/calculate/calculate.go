// Package calculate contains calculators for envelope punch positions.
package calculate

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}

	return AbsInt64(a)
}

// AbsInt64 returns the absolute value of an integer.
func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Math.gcd= function(a, b){
// 	if(b) return Math.gcd(b, a%b);
// 	return Math.abs(a);
// }

// A Rational number is expressed as the fraction p/q of two integers:
// r = p/q = (d*i+n)/d.
type Rational struct {
	i int64 // integer
	n int64 // fraction numerator
	d int64 // fraction denominator
}

func (r Rational) String() string {
	const base10 = 10

	var s string
	if r.i != 0 {
		s += strconv.FormatInt(r.i, base10)
	}
	if r.n != 0 {
		if r.i != 0 {
			s += " + "
		}
		if r.d < 0 {
			r.n *= -1
			r.d *= -1
		}
		s += strconv.FormatInt(r.n, base10) + "/" + strconv.FormatInt(r.d, base10)
	}
	if len(s) == 0 {
		s += "0"
	}
	return s
}

// ParseDecimal parses a decimal string into a rational number.
func ParseDecimal(s string) (r Rational, err error) {
	const (
		bitSize int = 64
		base10  int = 10
	)

	sign := int64(1)

	if strings.HasPrefix(s, "-") {
		sign = -1
	}

	p := strings.IndexByte(s, '.')

	if p < 0 {
		p = len(s)
	}

	if i := s[:p]; len(i) > 0 {
		if i != "+" && i != "-" {
			r.i, err = strconv.ParseInt(i, 10, 64)
			if err != nil {
				return Rational{}, err
			}
		} else {
			r.i = 0
		}
	}
	if p >= len(s) {
		p = len(s) - 1
	}

	if f := s[p+1:]; len(f) > 0 {
		n, err := strconv.ParseUint(f, base10, bitSize)
		if err != nil {
			return Rational{}, err
		}

		d := math.Pow10(len(f))

		if math.Log2(d) > float64(bitSize-1) {
			err = fmt.Errorf("ParseDecimal: parsing %q: value out of range", f)
			return Rational{}, err
		}

		r.n = int64(n)

		if r.n != 0 {
			r.d = int64(d)
		}

		if g := gcd(r.n, r.d); g != 0 {
			r.n /= g
			r.d /= g
		}

		if r.i == 0 {
			r.n *= sign
		}
	}

	return r, nil
}

// Math.fraction= function(n, prec, up){
// 		var s= String(n),
// 	p= s.indexOf('.');
// 	if(p== -1) return s;

// 	var i= Math.floor(n) || '',
// 	dec= s.substring(p),
// 	m= prec || Math.pow(10, dec.length-1),
// 	num= up=== 1? Math.ceil(dec*m): Math.round(dec*m),
// 	den= m,
// 	g= Math.gcd(num, den);

// 	if(den/g==1) return String(i+(num/g));

// 	if(i) i= i+' ';
// 	return i+ String(num/g)+'/'+String(den/g);
// }

// CalculateEnvelope calculates the paper size and punch location for an envelope.
func CalculateEnvelope(length, width float64, isLoose bool, boardMini bool) (float64, float64, error) {
	const (
		marginMetric            float64 = 1.1
		marginMetricLoose       float64 = 1.5
		marginImperial          float64 = 0.4375
		marginImperialLoose     float64 = 0.625
		marginMiniMetric        float64 = 0.67
		marginMiniMetricLoose   float64 = 0.914
		marginMiniImperial      float64 = 0.25
		marginMiniImperialLoose float64 = 0.34
		metricUnit              string  = "cm"
		imperialUnit            string  = "inch"
	)

	margin := 1.1 // TODO: figure out how to incorporate mini board values
	if isLoose {
		margin += 0.4
	}

	//   var isNotThick = $("#cardsizeb").is(":checked") || $("#cardsizec").is(":checked")
	//   var margin = isMini
	// 	? (isNotThick
	// 			 ? (imperial ? marginMiniImp : marginMiniMet )
	// 			 : (imperial ? marginMiniImpThick : marginMiniMetThick))
	// 	: (isNotThick
	// 	  ? (imperial ? marginImp : marginMet )
	// 	  : (imperial ? marginImpThick : marginMetThick));

	//   var boxMode = $("#cardsizec").is(":checked");
	//   $("#cardsizeSettingHeight").css("display", boxMode ? "" : "none");
	//   $("#cardsizeTrPunchPoint2").css("display", boxMode ? "" : "none");
	//   $("#punchpoint1").css("display", boxMode ? "" : "none");
	//   var slLength = GetNumericValue($("#cardsizeLength").val());
	//   var slWidth = GetNumericValue($("#cardsizeWidth").val());
	//   var slHeight = GetNumericValue($("#cardsizeHeight").val());

	const distMultiplier float64 = 0.707106781187 // 1/sqrt(2)
	dist1 := length * distMultiplier
	dist2 := width * distMultiplier

	//   var dist1 = slLength * Math.sqrt(0.5);
	//   var dist2 = slWidth * Math.sqrt(0.5);
	//   var dist3 = slHeight * Math.sqrt(0.5);

	//   var canvas = document.getElementById('cardsizeResultCanvas');
	//   var context = canvas.getContext ? canvas.getContext('2d') : null;
	//   var generalfillStyle = "#38761D";
	//   var boxFlapFillStyle = "#739F61";
	//   if (dist1 == 0 || dist2 == 0 || (dist3==0 && boxMode) )
	//   {
	// 	  $("#cardsizeResPunchPoint").text("");
	// 	  $("#cardsizeResPunchPoint2").text("");
	// 	  $("#cardsizeResPaperSize").text("");
	// 	context.clearRect(0, 0, canvas.height, canvas.height);
	//   }
	//   else
	//   {

	paper := dist1 + dist2 + 2*margin

	// 		 var paperSize = boxMode ? dist1 + dist2 + 2 * (dist3 + margin) : dist1 + dist2 + 2 * margin;
	// 		 var drawFactor = paperSize==0 ? 10 : (500.0 / paperSize);
	// 	 var drawWidth = paperSize * drawFactor;
	// 	 canvas.height = drawWidth;
	// 	 canvas.width = drawWidth;
	// 	 var makeCorner = function(cornerX, cornerY, vanDX, naarDX, vanDY, naarDY){
	// 			context.beginPath();
	// 			context.strokeStyle = "#888888";
	// 			context.setLineDash([5,2]);
	// 			var extraLine = margin;
	// 			context.moveTo((cornerX+vanDX*extraLine) * drawFactor, (cornerY+vanDY*extraLine) * drawFactor);
	// 			context.lineTo(cornerX * drawFactor, cornerY * drawFactor);
	// 			context.lineTo((cornerX+naarDX*extraLine) * drawFactor, (cornerY+naarDY*extraLine) * drawFactor);
	// 			context.stroke();
	// 			context.closePath();
	// 			context.fillStyle = "#e0e0e0";
	// 			context.fill();
	// 	  };
	// 	  if (boxMode)
	// 	  {
	// 		  if ( dist2 > dist1)
	// 		  {
	// 			$("#cardsizeResPunchPoint").html(FormatDistance(dist1 + margin,imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist1,imperial) + ")</span>");
	// 			$("#cardsizeResPunchPoint2").html(FormatDistance(dist1 + margin + 2 * dist3,imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist1,imperial) + " + 2 * " + FormatDistance(dist3,imperial) + ")</span>");

	// 			  context.clearRect(0, 0, drawWidth, drawWidth);
	// 			  context.fillStyle = generalfillStyle;
	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3) * drawFactor, (margin + dist3 + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist2) * drawFactor, (margin + dist3 + dist1 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist1 + dist2) * drawFactor, (margin + dist3 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist1) * drawFactor, (margin + dist3) * drawFactor);
	// 			  context.lineTo((margin + dist3) * drawFactor, (margin + dist3 + dist1) * drawFactor);
	// 			  context.closePath();
	// 			  context.fill();

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3) * drawFactor, (margin + dist3 + dist1) * drawFactor);
	// 			  context.lineTo((margin) * drawFactor, (margin + dist3 + dist3 + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist2) * drawFactor, (paperSize - margin) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist2) * drawFactor, (margin + dist3 + dist1 + dist2) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 				  makeCorner(margin, margin + dist3 + dist3 + dist1, -1, -1, -1, 1);
	// 				  makeCorner(margin + dist2, paperSize - margin, -1, 1, 1, 1);

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3 + dist2) * drawFactor, (margin + dist3 + dist1 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist3 + dist2) * drawFactor, (paperSize - margin) * drawFactor);
	// 			  context.lineTo((paperSize - margin) * drawFactor, (margin + dist3 + dist3 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist1 + dist2) * drawFactor, (margin + dist3 + dist2) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 				  makeCorner(margin + dist3 + dist3 + dist2, paperSize - margin, -1, 1, 1, 1);
	// 				  makeCorner(paperSize - margin, margin + dist3 + dist3 + dist2, 1, 1, -1, 1);

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3 + dist1 + dist2) * drawFactor, (margin + dist3 + dist2) * drawFactor);
	// 			  context.lineTo((paperSize - margin) * drawFactor, (margin + dist2) * drawFactor);
	// 			  context.lineTo((paperSize - margin - dist2) * drawFactor, (margin) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist1) * drawFactor, (margin + dist3) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 				  makeCorner(paperSize - margin, margin + dist2, 1, 1, -1, 1);
	// 				  makeCorner(paperSize - margin - dist2, margin, -1, 1, -1, -1);

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3 + dist1) * drawFactor, (margin + dist3) * drawFactor);
	// 			  context.lineTo((margin + dist1) * drawFactor, (margin) * drawFactor);
	// 			  context.lineTo((margin) * drawFactor, (margin + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist3) * drawFactor, (margin + dist3 + dist1) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 				  makeCorner(margin + dist1, margin, -1, 1, -1, -1);
	// 				  makeCorner(margin, margin + dist1, -1, -1, -1, 1);
	// 		  }
	// 		  else
	// 		  {
	// 			$("#cardsizeResPunchPoint").html(FormatDistance(dist2 + margin,imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist2,imperial) + ")</span>");
	// 			$("#cardsizeResPunchPoint2").html(FormatDistance(dist2 + margin + 2 * dist3,imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist2,imperial) + " + 2 * " + FormatDistance(dist3,imperial) + ")</span>");

	// 			  context.clearRect(0, 0, drawWidth, drawWidth);
	// 			  context.fillStyle = generalfillStyle;
	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3) * drawFactor, (margin + dist3 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist1) * drawFactor, (margin + dist3 + dist2 + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist2 + dist1) * drawFactor, (margin + dist3 + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist2) * drawFactor, (margin + dist3) * drawFactor);
	// 			  context.lineTo((margin + dist3) * drawFactor, (margin + dist3 + dist2) * drawFactor);
	// 			  context.closePath();
	// 			  context.fill();

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3) * drawFactor, (margin + dist3 + dist2) * drawFactor);
	// 			  context.lineTo((margin) * drawFactor, (margin + dist3 + dist3 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist1) * drawFactor, (paperSize - margin) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist1) * drawFactor, (margin + dist3 + dist2 + dist1) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 				makeCorner(margin, margin + dist3 + dist3 + dist2, -1, -1, -1, 1);
	// 			makeCorner(margin + dist1, paperSize - margin, -1, 1, 1, 1);

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3 + dist1) * drawFactor, (margin + dist3 + dist2 + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist3 + dist1) * drawFactor, (paperSize - margin) * drawFactor);
	// 			  context.lineTo((paperSize - margin) * drawFactor, (margin + dist3 + dist3 + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist2 + dist1) * drawFactor, (margin + dist3 + dist1) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 			makeCorner(margin + dist3 + dist3 + dist1, paperSize - margin, -1, 1, 1, 1);
	// 			makeCorner(paperSize - margin, margin + dist3 + dist3 + dist1, 1, 1, -1, 1);

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3 + dist2 + dist1) * drawFactor, (margin + dist3 + dist1) * drawFactor);
	// 			  context.lineTo((paperSize - margin) * drawFactor, (margin + dist1) * drawFactor);
	// 			  context.lineTo((paperSize - margin - dist1) * drawFactor, (margin) * drawFactor);
	// 			  context.lineTo((margin + dist3 + dist2) * drawFactor, (margin + dist3) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 			makeCorner(paperSize - margin, margin + dist1, 1, 1, -1, 1);
	// 			makeCorner(paperSize - margin - dist1, margin, -1, 1, -1, -1);

	// 			  context.beginPath();
	// 			  context.moveTo((margin + dist3 + dist2) * drawFactor, (margin + dist3) * drawFactor);
	// 			  context.lineTo((margin + dist2) * drawFactor, (margin) * drawFactor);
	// 			  context.lineTo((margin) * drawFactor, (margin + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist3) * drawFactor, (margin + dist3 + dist2) * drawFactor);
	// 			  context.closePath();
	// 			  context.fillStyle = boxFlapFillStyle;
	// 			  context.fill();

	// 			makeCorner(margin + dist2, margin, -1, 1, -1, -1);
	// 			makeCorner(margin, margin + dist2, -1, -1, -1, 1);
	// 		  }
	// 		  $("#cardsizeResPaperSize").html(FormatDistance(dist1 + dist2 + 2 * (dist3 + margin),imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist1,imperial) + " + " + FormatDistance(dist2,imperial) + " + 2 * " + FormatDistance(dist3,imperial) + " + " + FormatDistance(margin,imperial) + ")</span>");
	// 	  }
	// 	  else
	// 	  {
	punch := margin + math.Min(dist1, dist2)

	// 		  if ( dist2 > dist1 )
	// 		  {
	// 			$("#cardsizeResPunchPoint").html(FormatDistance(dist1 + margin,imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist1,imperial) + ")</span>");

	// 			  context.clearRect(0, 0, drawWidth, drawWidth);
	// 			  context.fillStyle = generalfillStyle;
	// 			  context.beginPath();
	// 			  context.moveTo(margin * drawFactor, (margin + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist2) * drawFactor, (margin + dist1 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist1 + dist2) * drawFactor, (margin + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist1) * drawFactor, (margin) * drawFactor);
	// 			  context.lineTo((margin) * drawFactor, (margin + dist1) * drawFactor);
	// 			  context.closePath();
	// 			  context.fill();

	// 		  makeCorner(margin, margin + dist1, -1, -1, -1, 1);
	// 				  makeCorner(margin + dist2, margin + dist1+dist2, -1, 1, 1, 1);
	// 				  makeCorner(margin + dist1 +dist2, margin + dist2, 1, 1, -1, 1);
	// 				  makeCorner(margin + dist1, margin, -1, 1, -1, -1);
	// 		  }
	// 		  else
	// 		  {
	// 			$("#cardsizeResPunchPoint").html(FormatDistance(dist2 + margin,imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist2,imperial) + ")</span>");

	// 			  context.clearRect(0, 0, drawWidth, drawWidth);
	// 			  context.fillStyle = generalfillStyle;
	// 			  context.beginPath();
	// 			  context.moveTo(margin * drawFactor, (margin + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist1) * drawFactor, (margin + dist1 + dist2) * drawFactor);
	// 			  context.lineTo((margin + dist1 + dist2) * drawFactor, (margin + dist1) * drawFactor);
	// 			  context.lineTo((margin + dist2) * drawFactor, (margin) * drawFactor);
	// 			  context.lineTo((margin) * drawFactor, (margin + dist2) * drawFactor);
	// 			  context.closePath();
	// 			  context.fill();

	// 				  makeCorner(margin, margin + dist2, -1, -1, -1, 1);
	// 				  makeCorner(margin + dist1, margin + dist1+dist2, -1, 1, 1, 1);
	// 				  makeCorner(margin + dist1 +dist2, margin + dist1, 1, 1, -1, 1);
	// 				  makeCorner(margin + dist2, margin, -1, 1, -1, -1);
	// 		  }
	// 		  $("#cardsizeResPaperSize").html(FormatDistance(paperSize,imperial) + unitCode + " <span class=\"calcdetails\">(= " + FormatDistance(margin,imperial) + " + " + FormatDistance(dist1,imperial) + " + " + FormatDistance(dist2,imperial) + " + " + FormatDistance(margin,imperial) + ")</span>");
	// 	  }
	//   }
	// }

	// $(document).on("ready", function ()
	// {
	//   InitMode($("#unitimp").is(":checked"));
	//   $("#cardsizea,#cardsizeb,#cardsizec,#unitmet,#unitimp").click(function(){
	// 		InitMode($("#unitimp").is(":checked"));
	// 		CalculateSizes();
	// 	});
	//   $("#boardtype").on("change",function(){
	// 		CalculateSizes();
	// 	});
	//   $("#cardsizeWidth,#cardsizeLength,#cardsizeHeight").on("change keyup focusout", CalculateSizes);

	//   $("#cardsizeWidth,#cardsizeLength,#cardsizeHeight").on("change focusout", function(){
	// 	   InitMode($("#unitimp").is(":checked"));
	//    });

	//   $("#cardsizeWidth,#cardsizeLength,#cardsizeHeight").on("focus", function ()
	// 	{
	// 	  var v = $(this).val();
	// 	  if( v==formatMet || v==formatImp )
	// 	  {
	// 		 $(this).select();
	// 	  }
	//    });
	//   CalculateSizes();
	// });
	return paper, punch, nil
}
