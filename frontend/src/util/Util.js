export const indicatorToText = (indicator) => {
    switch (parseInt(indicator)) {
        case 0:
            return "Completely Disagree"

        case 1:
        case 2:
            return "Disagree"

        case 3:
        case 4:
            return "Slightly Disagree"
        
        case 5:
            return "Neutral"

        case 6:
        case 7:
            return "Slightly Agree"

        case 8:
        case 9:
            return "Agree"
            
        case 10:
            return "Completely Agree"
    
        default:
            return "No Opinion"
    }
}