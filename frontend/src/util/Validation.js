// Validations and messages presume submition buttons are disabled if either topic or detail is blank

export const validateFields = (topic, indicator, detail) => {
    return {
        topic: validateTitle(topic),
        indicator: indicator >= 0 && indicator <= 10,
        detail: validateDetail(detail),
    }
}

export const validateTitle = (topic) => {
    return topic.length <= 140
}

export const validateDetail = (detail) => {
    return detail.length <= 1250
}

export const getTitleValidationMessage = () => {
    return "must be less than 140 characters"
}

export const getIndicatorValidationMessage = () => {
    return "must be between 0 and 10 included"
}

export const getDetailValidationMessage = () => {
    return "must be less than 1250 characters"
}