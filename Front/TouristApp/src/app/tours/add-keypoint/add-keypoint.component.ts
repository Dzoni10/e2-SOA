// add-keypoint-modal.component.ts
import { Component, ElementRef, Inject, ViewChild, AfterViewInit, OnDestroy } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';

declare var L: any; // Leaflet

interface SelectedImage {
  file: File;
  preview: string;
}

@Component({
  selector: 'app-add-keypoint',
  templateUrl: './add-keypoint.component.html',
  styleUrls: ['./add-keypoint.component.css']
  
})
export class AddKeypointComponent implements AfterViewInit, OnDestroy {
  @ViewChild('mapContainer', { static: false }) mapContainer!: ElementRef;

  keypointForm: FormGroup;
  selectedImages: SelectedImage[] = [];
  gettingLocation = false;
  isSubmitting = false;

  private map: any;
  private marker: any;

  constructor(
    private fb: FormBuilder,
    private dialogRef: MatDialogRef<AddKeypointComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {
    this.keypointForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      latitude: ['', [Validators.required, Validators.min(-90), Validators.max(90)]],
      longitude: ['', [Validators.required, Validators.min(-180), Validators.max(180)]]
    });

    // Subscribe to coordinate changes to update map
    this.keypointForm.get('latitude')?.valueChanges.subscribe(() => this.updateMapFromForm());
    this.keypointForm.get('longitude')?.valueChanges.subscribe(() => this.updateMapFromForm());
  }

  ngAfterViewInit() {
    this.initializeMap();
  }

  ngOnDestroy() {
    if (this.map) {
      this.map.remove();
    }
    // Clean up image previews
    this.selectedImages.forEach(img => {
      URL.revokeObjectURL(img.preview);
    });
  }

  private initializeMap() {
    // Default to Belgrade coordinates
    const defaultLat = 44.8176;
    const defaultLng = 20.4633;

    this.map = L.map(this.mapContainer.nativeElement).setView([defaultLat, defaultLng], 13);

    // Add OpenStreetMap tiles
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '© OpenStreetMap contributors'
    }).addTo(this.map);

    // Add click handler
    this.map.on('click', (e: any) => {
      this.onMapClick(e);
    });
  }

  private onMapClick(e: any) {
    const lat = e.latlng.lat;
    const lng = e.latlng.lng;

    // Update form
    this.keypointForm.patchValue({
      latitude: lat.toFixed(6),
      longitude: lng.toFixed(6)
    });

    // Update marker
    this.updateMarker(lat, lng);
  }

  private updateMapFromForm() {
    const lat = this.keypointForm.get('latitude')?.value;
    const lng = this.keypointForm.get('longitude')?.value;

    if (lat && lng && !isNaN(lat) && !isNaN(lng)) {
      const numLat = parseFloat(lat);
      const numLng = parseFloat(lng);
      
      if (numLat >= -90 && numLat <= 90 && numLng >= -180 && numLng <= 180) {
        this.map.setView([numLat, numLng], this.map.getZoom());
        this.updateMarker(numLat, numLng);
      }
    }
  }

  private updateMarker(lat: number, lng: number) {
    if (this.marker) {
      this.marker.setLatLng([lat, lng]);
    } else {
      this.marker = L.marker([lat, lng], {
        draggable: true
      }).addTo(this.map);

      // Handle marker drag
      this.marker.on('dragend', (e: any) => {
        const position = e.target.getLatLng();
        this.keypointForm.patchValue({
          latitude: position.lat.toFixed(6),
          longitude: position.lng.toFixed(6)
        });
      });
    }
  }

  getCurrentLocation() {
    if (!navigator.geolocation) {
      alert('Geolocation is not supported by this browser.');
      return;
    }

    this.gettingLocation = true;

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const lat = position.coords.latitude;
        const lng = position.coords.longitude;

        this.keypointForm.patchValue({
          latitude: lat.toFixed(6),
          longitude: lng.toFixed(6)
        });

        this.map.setView([lat, lng], 15);
        this.updateMarker(lat, lng);
        
        this.gettingLocation = false;
      },
      (error) => {
        console.error('Error getting location:', error);
        alert('Unable to get current location. Please set location manually.');
        this.gettingLocation = false;
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 300000
      }
    );
  }

  onFilesSelected(event: any) {
    const files = Array.from(event.target.files) as File[];
    
    files.forEach(file => {
      if (file.type.startsWith('image/')) {
        const preview = URL.createObjectURL(file);
        this.selectedImages.push({ file, preview });
      }
    });

    // Clear input
    event.target.value = '';
  }

  removeImage(index: number) {
    const image = this.selectedImages[index];
    URL.revokeObjectURL(image.preview);
    this.selectedImages.splice(index, 1);
  }

  async addKeypoint() {
    if (!this.keypointForm.valid) {
      return;
    }

    this.isSubmitting = true;

    try {
      const formData = new FormData();
      const formValue = this.keypointForm.value;

      // Add form data
      formData.append('name', formValue.name);
      formData.append('description', formValue.description);
      formData.append('latitude', formValue.latitude);
      formData.append('longitude', formValue.longitude);

      // Add images
      this.selectedImages.forEach(image => {
        formData.append('images', image.file);
      });

      // Create temporary keypoint object for preview
      const keypointData = {
        name: formValue.name,
        description: formValue.description,
        latitude: parseFloat(formValue.latitude),
        longitude: parseFloat(formValue.longitude),
        images: this.selectedImages.map(img => img.preview), // Use preview URLs for now
        formData: formData // Include form data for actual submission
      };

      // Return the keypoint data
      this.dialogRef.close(keypointData);

    } catch (error) {
      console.error('Error creating keypoint:', error);
      alert('Error creating keypoint. Please try again.');
    } finally {
      this.isSubmitting = false;
    }
  }
}
