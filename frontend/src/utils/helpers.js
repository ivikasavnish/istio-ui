import yaml from 'js-yaml';

export const objectToYaml = (obj) => {
  try {
    return yaml.dump(obj, { indent: 2, lineWidth: -1 });
  } catch (error) {
    console.error('Error converting to YAML:', error);
    return '';
  }
};

export const yamlToObject = (yamlStr) => {
  try {
    return yaml.load(yamlStr);
  } catch (error) {
    console.error('Error parsing YAML:', error);
    throw error;
  }
};

export const formatDate = (date) => {
  if (!date) return 'N/A';
  return new Date(date).toLocaleString();
};

export const truncate = (str, maxLength = 50) => {
  if (!str) return '';
  return str.length > maxLength ? str.substring(0, maxLength) + '...' : str;
};
