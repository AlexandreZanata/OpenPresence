use image::{DynamicImage, RgbImage};

const RECOGNITION_SIZE: u32 = 112;

/// Affine warp simplified: scale to 112×112 RGB, normalize to [-1, 1] (mean/std 0.5).
pub fn preprocess_for_recognition(img: &DynamicImage, _landmarks: &[[f32; 2]; 5]) -> Vec<f32> {
    let resized = img.resize_exact(RECOGNITION_SIZE, RECOGNITION_SIZE, image::imageops::FilterType::Triangle);
    let rgb = resized.to_rgb8();
    rgb_normalized(&rgb)
}

fn rgb_normalized(rgb: &RgbImage) -> Vec<f32> {
    let mut out = Vec::with_capacity((rgb.width() * rgb.height() * 3) as usize);
    for pixel in rgb.pixels() {
        for c in 0..3 {
            out.push((pixel[c] as f32 / 255.0 - 0.5) / 0.5);
        }
    }
    out
}

#[cfg(test)]
mod tests {
    use super::*;
    use image::RgbImage;

    #[test]
    fn preprocess_recognition_output_size() {
        let img = DynamicImage::ImageRgb8(RgbImage::new(200, 200));
        let landmarks = [[0.0_f32; 2]; 5];
        let data = preprocess_for_recognition(&img, &landmarks);
        assert_eq!(data.len(), 112 * 112 * 3);
    }
}
