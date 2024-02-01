
export const validateFields = (value: any) => {
    const allowedFields = ['ip', 'port', 'peerId'];
  
    if (!value || typeof value !== 'object') {
      throw new Error('Invalid request body');
    }
  
    const keys = Object.keys(value);
    const invalidKeys = keys.filter((key) => !allowedFields.includes(key));
  
    if (invalidKeys.length > 0) {
      throw new Error(`Invalid field(s): ${invalidKeys.join(', ')}`);
    }
  
    return true;
};
