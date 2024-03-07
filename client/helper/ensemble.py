import torch
import os

def average_state_dicts(state_dict_files):
    # Load the first state dictionary
    model_state_dict = torch.load(state_dict_files[0])
    
    # Initialize parameters
    num_state_dicts = len(state_dict_files)

    # Iterate over keys in the state dictionary
    for key in model_state_dict.keys():
        if torch.is_tensor(model_state_dict[key]):  # Check if the value is a tensor
            if model_state_dict[key].dtype == torch.float32:  # Check if the tensor has float32 dtype
                model_state_dict[key] = model_state_dict[key].clone().float()  # Convert to float if necessary
            elif model_state_dict[key].dtype == torch.int64:  # Check if the tensor has int64 dtype
                model_state_dict[key] = model_state_dict[key].clone().float()  # Convert to float if necessary
            else:
                raise ValueError("Unsupported data type: {}".format(model_state_dict[key].dtype))

    # Iterate over remaining state dictionaries and add their parameters
    for i in range(1, num_state_dicts):
        current_state_dict = torch.load(state_dict_files[i])
        for key in model_state_dict.keys():
            model_state_dict[key] += current_state_dict[key]

    # Calculate the average
    for key in model_state_dict.keys():
        model_state_dict[key] /= num_state_dicts
    
    return model_state_dict

if __name__ == "__main__":
    state_dict_folder = 'weights'  # Folder containing model state dictionary files
    state_dict_files = [os.path.join(state_dict_folder, f) for f in os.listdir(state_dict_folder) if f.endswith('.pth')]
    final_state_dict = average_state_dicts(state_dict_files)
    torch.save(final_state_dict, 'averaged_model_state_dict.pth')
