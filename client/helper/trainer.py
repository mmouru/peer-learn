""" Training iterations for resnet 18 models """

import torch
import torchvision
from torchvision import datasets
import torchvision.transforms as transforms
import torch.nn as nn
import torch.optim as optim
import glob
import os

## Confirm the device
#device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
#print(device)
device = torch.device("cpu")

## We already need to have normalized data that is passed to clients

## use resnet18 model with pretrained weights
model = torchvision.models.resnet18(weights="ResNet18_Weights.DEFAULT")

criterion = nn.CrossEntropyLoss()
optimizer = optim.SGD(model.parameters(), lr=0.001, momentum=0.9)
current_directory = os.getcwd()
print("Current working directory:", current_directory)
subfolders = glob.glob('data/all_training_data/training_set/')

train_data_path = "data/all_training_data/training_set/"

transform = transforms.Compose([
    # Add your transformations here, e.g., resizing, normalization, data augmentation
    transforms.ToTensor(),
])

imgfolder = datasets.ImageFolder(root=train_data_path, transform=transform)
n_classes = len(imgfolder.classes)

num_ftrs = model.fc.in_features
model.fc = nn.Linear(num_ftrs, n_classes)
model.to(device)

trainloader = torch.utils.data.DataLoader(imgfolder, batch_size=256, shuffle=True, num_workers=2)

total_size = 0

# Iterate through the dataset

print(len(imgfolder.imgs))

## training function

def train_for_epochs(epochs=10):
    
    for epoch in range(epochs):  # loop over the dataset multiple times
        running_loss = 0.0
        for data in trainloader:# get the inputs; data is a list of [inputs, labels]
            inputs, labels = data
            inputs = inputs.to(device)
            labels = labels.to(device)

            # zero the parameter gradients
            optimizer.zero_grad()

            # forward + backward + optimize
            outputs = model(inputs)
            
            _, predicted = torch.max(outputs, 1)  # Get predicted labels
            
            loss = criterion(outputs, labels)
            loss.backward()
            optimizer.step()

            # print statistics
            running_loss += loss.item()
        print("epoch ", epoch, "running loss:", running_loss)
    
    torch.save(model.state_dict(), './data/model_state_dict.pth')
    

if __name__ == "__main__":
    train_for_epochs(8)