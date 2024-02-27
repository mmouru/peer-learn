import { useState } from 'react';

export const showFileSize = (bytes) => {
    const gigabyte = 1024 * 1024 * 1024;
    const megabyte = 1024 * 1024;
    if (bytes >= gigabyte) {
        return (bytes / gigabyte).toFixed(2) + ' GB';
    } else if (bytes >= megabyte) {
        return (bytes / megabyte).toFixed(2) + ' MB';
    } else {
        return bytes + ' bytes';
    }
}

export const checkForZipType = (file) => {
    if (file.name && file.type != "application/x-zip-compressed") {
        return true
    }
    return false
};


export const FileSelection = (props) => {
    const [selectedFile, setSelectedFiles] = useState({name: "", size: 0, type: ""});

    const handleFileChange = (event) => {
        const files = event.target.files;
        console.log(files)
        setSelectedFiles(files[0]);
      };

    const showFileSize = (bytes) => {
        const gigabyte = 1024 * 1024 * 1024;
        const megabyte = 1024 * 1024;
        if (bytes >= gigabyte) {
            return (bytes / gigabyte).toFixed(2) + ' GB';
        } else if (bytes >= megabyte) {
            return (bytes / megabyte).toFixed(2) + ' MB';
        } else {
            return bytes + ' bytes';
        }
    }

    return (
        <div>
        <p> File drop </p>
        <input type="file" onChange={handleFileChange}/>
        <table>
          <thead>
            <tr>
              <th>File Name</th>
              <th>File Size</th>
            </tr>
          </thead>
          <tbody>
            <tr>
                <td>
                    {selectedFile.name}
                </td>
                <td>
                    {showFileSize(selectedFile.size)}
                </td>
            </tr>
          </tbody>
        </table>
        {checkForZipType(selectedFile) ? <p style={{color: "red"}}>Requires zip file!</p> : ""}
      </div>
    )
}